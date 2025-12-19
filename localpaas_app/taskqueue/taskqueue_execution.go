package taskqueue

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
)

const (
	taskInfoCacheExp    = 24 * time.Hour
	taskRetryMaxBackoff = 24 * time.Hour
)

func (q *taskQueue) loadTask(
	ctx context.Context,
	db database.IDB,
	taskID string,
) (*entity.Task, error) {
	task, err := q.taskRepo.GetByID(ctx, db, "", taskID,
		bunex.SelectWhereIn("task.status IN (?)", base.TaskStatusNotStarted, base.TaskStatusFailed),
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
	executor func(context.Context, database.Tx, *entity.Task) error,
) (rescheduleAt time.Time, err error) {
	var task, nextTask *entity.Task
	err = transaction.Execute(ctx, q.db, func(db database.Tx) error {
		task, err = q.loadTask(ctx, db, taskID)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if task == nil {
			return nil
		}
		if task.NextTaskID != "" {
			nextTask, err = q.loadTask(ctx, db, task.NextTaskID)
			if err != nil {
				return apperrors.Wrap(err)
			}
		}

		// Put task to in-progress state
		task.StartedAt = timeutil.NowUTC()
		if task.Status == base.TaskStatusFailed {
			task.Retry++
		}

		// Mark the task as `in-progress` by inserting a new record in to cache
		taskInfo := &cacheentity.TaskInfo{
			ID:        task.ID,
			Status:    base.TaskStatusInProgress,
			StartedAt: task.StartedAt,
		}
		err = q.cacheTaskInfoRepo.Set(ctx, task.ID, taskInfo, taskInfoCacheExp)
		if err != nil {
			return apperrors.Wrap(err)
		}

		defer func() {
			_ = q.cacheTaskInfoRepo.Del(ctx, task.ID)
			err = q.taskRepo.UpdateMulti(ctx, db, gofn.ToSliceSkippingNil(task, nextTask))
		}()

		err = executor(ctx, db, task)
		timeNow := timeutil.NowUTC()
		if err != nil {
			nextTask = nil
			task.Status = base.TaskStatusFailed
			task.RetryAt = task.EndedAt.Add(calcExpBackoffRetry(task))
			_ = task.AddRun(&entity.TaskRun{
				StartedAt: task.StartedAt,
				EndedAt:   timeNow,
				Error:     err.Error(),
			})
			return apperrors.Wrap(err)
		}

		task.Status = base.TaskStatusDone
		task.EndedAt = timeNow
		if nextTask != nil {
			nextTask.RunAt = timeNow
		}
		return nil
	})
	if err != nil {
		if task != nil && task.Status == base.TaskStatusFailed {
			rescheduleAt = task.RetryAt
		}
		return rescheduleAt, apperrors.Wrap(err)
	}

	task = nextTask
	if task != nil {
		err = q.server.ScheduleNextTask(task, time.Time{})
		if err != nil {
			if task != nil && task.Status == base.TaskStatusFailed {
				rescheduleAt = task.RetryAt
			}
			return rescheduleAt, apperrors.Wrap(err)
		}
	}

	return rescheduleAt, nil
}

func calcExpBackoffRetry(task *entity.Task) time.Duration {
	randDur := time.Duration(rand.Int31n(1000)) * time.Millisecond //nolint:mnd,gosec
	if task.RetryDelaySecs == 0 {
		return randDur
	}
	exp := 1.0
	if task.Retry > 0 {
		exp = math.Pow(2, float64(task.Retry)) //nolint:mnd
	}
	return min(time.Duration(exp*float64(task.RetryDelaySecs))*time.Second+randDur, taskRetryMaxBackoff)
}
