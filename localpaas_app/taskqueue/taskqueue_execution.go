package taskqueue

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/localpaas/localpaas/localpaas_app/pkg/redishelper"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
)

const (
	taskDefaultTimeout      = 3 * time.Hour
	taskInfoCacheExp        = 24 * time.Hour
	taskRetryMaxBackoff     = 24 * time.Hour
	taskControlCheckTimeout = 5 * time.Second
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
		if errors.Is(err, apperrors.ErrNotFound) { // task not found, it's not error
			return nil, nil
		}
		return nil, apperrors.Wrap(err)
	}
	if task.JobID != "" {
		// Task's job is not active
		if task.Job == nil || !task.Job.IsActive() {
			return nil, nil
		}
	}
	// Task not allow retrying
	if task.Status == base.TaskStatusFailed && !task.CanRetry() {
		return nil, nil
	}
	return task, nil
}

func (q *taskQueue) executeTask(
	ctx context.Context,
	taskID string,
	_ string,
	executorFunc func(context.Context, database.Tx, *TaskExecData) error,
) (rescheduleAt time.Time, err error) {
	var task *entity.Task
	err = transaction.Execute(ctx, q.db, func(db database.Tx) error {
		task, err = q.loadTask(ctx, db, taskID)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if task == nil {
			return nil
		}

		taskData := &TaskExecData{
			Task: task,
		}
		taskTimeout := task.Config.Timeout.ToDuration()
		ctx, cancel := context.WithTimeout(ctx, gofn.Coalesce(taskTimeout, taskDefaultTimeout))
		defer cancel()

		// Check for commands from other processes
		go q.taskControlCheck(ctx, taskData)

		// Put task to in-progress state
		task.StartedAt = timeutil.NowUTC()
		if task.Status == base.TaskStatusFailed {
			task.Config.Retry++
		}

		// Mark the task as `in-progress` by inserting a new record in to redis
		taskInfo := &cacheentity.TaskInfo{
			ID:        task.ID,
			Status:    base.TaskStatusInProgress,
			StartedAt: task.StartedAt,
		}
		err = q.taskInfoRepo.Set(ctx, task.ID, taskInfo, taskInfoCacheExp)
		if err != nil {
			return apperrors.Wrap(err)
		}

		defer func() {
			if err == nil {
				if r := recover(); r != nil { // recover from panic
					err = apperrors.NewPanic(fmt.Sprintf("%v", r))
				}
			}
			_ = q.taskInfoRepo.Del(ctx, task.ID)
			err = q.taskRepo.Update(ctx, db, task)
		}()

		err = executorFunc(ctx, db, taskData)
		timeNow := timeutil.NowUTC()
		if err != nil {
			task.Status = base.TaskStatusFailed
			task.RetryAt = timeNow.Add(calcExpBackoffRetry(task))
			_ = task.AddRun(&entity.TaskRun{
				StartedAt: task.StartedAt,
				EndedAt:   timeNow,
				Error:     err.Error(),
			})
			return apperrors.Wrap(err)
		}

		task.EndedAt = timeNow
		task.Status = gofn.If(taskData.Canceled, base.TaskStatusCanceled, base.TaskStatusDone)
		return nil
	})
	if err != nil {
		if task != nil && task.Status == base.TaskStatusFailed {
			rescheduleAt = task.RetryAt
		}
		return rescheduleAt, apperrors.Wrap(err)
	}

	return rescheduleAt, nil
}

func (q *taskQueue) taskControlCheck(
	ctx context.Context,
	taskData *TaskExecData,
) {
	key := fmt.Sprintf("task:%s:ctrl", taskData.Task.ID)
	defer func() {
		_ = recover()
		_ = redishelper.Del(ctx, q.redisClient, key)
	}()

	for {
		if taskData.Done || taskData.Canceled {
			return
		}
		select {
		case <-ctx.Done(): // context is done, returns
			return
		default:
		}

		taskControls, err := redishelper.BLPop(ctx, q.redisClient, []string{key}, taskControlCheckTimeout,
			redishelper.JSONValueCreator[*cacheentity.TaskControl])
		if err != nil || len(taskControls) == 0 {
			continue
		}
		cmd := taskControls[key].Cmd
		if taskData.onCommand != nil {
			taskData.onCommand(cmd)
		}
		if taskData.Uncancelable && cmd == base.TaskCommandCancel {
			taskData.Canceled = true
			return
		}
	}
}

func calcExpBackoffRetry(task *entity.Task) time.Duration {
	randDur := time.Duration(rand.Int31n(1000)) * time.Millisecond //nolint:mnd,gosec
	delay := task.Config.RetryDelay.ToDuration()
	if delay == 0 {
		return randDur
	}
	exp := 1.0
	if task.Config.Retry > 0 {
		exp = math.Pow(2, float64(task.Config.Retry)) //nolint:mnd
	}
	return min(time.Duration(exp*float64(delay))+randDur, taskRetryMaxBackoff)
}
