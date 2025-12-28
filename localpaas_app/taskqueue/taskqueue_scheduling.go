package taskqueue

import (
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
)

const (
	missedTaskPeriod = 5 * time.Minute
)

func (q *taskQueue) doCreateTasks(
	ctx context.Context,
) error {
	err := transaction.Execute(ctx, q.db, func(db database.Tx) error {
		_, err := q.createTasks(ctx, db, nil, q.config.TaskQueue.TaskCreateInterval)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (q *taskQueue) doScheduleTasks(
	ctx context.Context,
) ([]*entity.Task, error) {
	timeNow := timeutil.NowUTC()
	tasks, _, err := q.taskRepo.List(ctx, q.db, "", nil,
		// Not-started tasks
		bunex.SelectWhereGroup(
			bunex.SelectWhere("task.status = ?", base.TaskStatusNotStarted),
			bunex.SelectWhere("(task.run_at IS NULL OR (task.run_at >= ? AND task.run_at < ?))",
				timeNow.Add(-missedTaskPeriod), timeNow.Add(q.config.TaskQueue.TaskCheckInterval)),
		),
		// Failed tasks need retry
		bunex.SelectWhereOrGroup(
			bunex.SelectWhere("task.status = ?", base.TaskStatusFailed),
			bunex.SelectWhere("task.max_retry > task.retry"),
			bunex.SelectWhere("task.retry_at IS NOT NULL"),
			bunex.SelectWhere("task.retry_at >= ?", timeNow.Add(-missedTaskPeriod)),
			bunex.SelectWhere("task.retry_at < ?", timeNow.Add(q.config.TaskQueue.TaskCheckInterval)),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(tasks) == 0 {
		return nil, nil
	}

	scheduleTasks := make([]*entity.Task, 0, len(tasks))
	for _, task := range tasks {
		if task.Status == base.TaskStatusNotStarted && task.RunAt.Before(timeNow) &&
			!q.shouldRunMissedTask(task, tasks, timeNow) {
			continue
		}
		scheduleTasks = append(scheduleTasks, task)
	}

	return scheduleTasks, nil
}

func (q *taskQueue) shouldRunMissedTask(
	missedTask *entity.Task,
	allTasks []*entity.Task,
	timeNow time.Time,
) bool {
	if missedTask.JobID == "" { // This is a solo task
		return true
	}
	for _, task := range allTasks {
		if task.JobID != missedTask.JobID || task.ID == missedTask.ID {
			continue
		}
		runAt := task.ShouldRunAt()
		if runAt.IsZero() {
			runAt = timeNow
		}
		if task.Status == base.TaskStatusNotStarted && runAt.Before(timeNow) && runAt.After(missedTask.RunAt) {
			return false
		}
		// The next run is near, so ignore the missed task?
		if runAt.Sub(timeNow) < timeNow.Sub(missedTask.RunAt) {
			return false
		}
	}
	return true
}
