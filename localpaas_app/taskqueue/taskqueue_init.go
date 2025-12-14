package taskqueue

import (
	"context"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/gocronqueue"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
)

func (q *taskQueue) Start(cfg *config.Config) error {
	// Initialize task queue worker if configured
	if cfg.RunMode == "worker" || cfg.RunMode == "embedded-worker" {
		go func() {
			q.logger.Infof("starting task queue server...")
			err := q.server.Start(gocronqueue.StartConfig{
				TaskCheckFunc: q.scanTasksForRunning,
				TaskMap: map[base.TaskType]gocronqueue.TaskProcessorFunc{
					base.TaskTypeTest: q.NewTestTaskProcessor(),
				},
			})
			if err != nil {
				q.logger.Errorf("failed to start task queue server: %v", err)
			}
		}()
	}
	if cfg.RunMode == "app" || cfg.RunMode == "embedded-worker" {
		q.logger.Infof("starting task queue client...")
	}

	return nil
}

func (q *taskQueue) Shutdown() error {
	q.logger.Info("stopping task queue ...")
	err := q.server.Shutdown()
	if err != nil {
		q.logger.Errorf("failed to start task queue server: %v", err)
		return apperrors.Wrap(err)
	}
	err = q.client.Close()
	if err != nil {
		q.logger.Errorf("failed to stop task queue client: %v", err)
		return apperrors.Wrap(err)
	}
	return nil
}

func (q *taskQueue) loadTask(
	ctx context.Context,
	db database.IDB,
	taskID string,
) (*entity.Task, error) {
	task, err := q.taskRepo.GetByID(ctx, db, "", taskID,
		bunex.SelectWhere("task.status IN (?)",
			bunex.InItems(base.TaskStatusNotStarted, base.TaskStatusFailed)),
		bunex.SelectFor("UPDATE OF task SKIP LOCKED"),
		bunex.SelectRelation("Job"),
	)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) { // task not found, it's no error
			return nil, nil
		}
		return nil, apperrors.Wrap(err)
	}
	if task.Job == nil || !task.Job.IsActive() {
		return nil, nil
	}
	// Task not allow retrying
	if task.Status == base.TaskStatusFailed && task.MaxRetry <= task.Retry {
		return nil, nil
	}
	return task, nil
}

func (q *taskQueue) runTask(
	ctx context.Context,
	taskID string,
	_ string,
	executor func(context.Context, database.IDB, *entity.Task) error,
) error {
	err := transaction.Execute(ctx, q.db, func(db database.Tx) error {
		task, err := q.loadTask(ctx, db, taskID)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if task == nil {
			return nil
		}

		// Put task to in-progress state
		task.StartedAt = timeutil.NowUTC()
		if task.Status == base.TaskStatusFailed {
			task.Retry++
		}

		// Mark the task as `in-progress` by inserting a new record
		// NOTE: we must use `q.db` to not bound the writing to this transaction
		inProgressTask := &entity.UpdatingTask{
			ID:        task.ID,
			StartedAt: task.StartedAt,
		}
		err = q.updatingTaskRepo.Upsert(ctx, q.db, inProgressTask,
			entity.UpdatingTaskUpsertingConflictCols, entity.UpdatingTaskUpsertingUpdateCols)
		if err != nil {
			return apperrors.Wrap(err)
		}

		defer func() {
			err = q.updatingTaskRepo.Delete(ctx, db, inProgressTask)
			err = q.taskRepo.Update(ctx, db, task)
		}()

		err = executor(ctx, db, task)
		task.Status = base.TaskStatusDone
		task.EndedAt = timeutil.NowUTC()

		if err != nil {
			task.Status = base.TaskStatusFailed
			_ = task.AddRun(&entity.TaskRun{
				StartedAt: task.StartedAt,
				EndedAt:   task.EndedAt,
				Error:     err.Error(),
			})
			return apperrors.Wrap(err)
		}

		return nil
	})
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
