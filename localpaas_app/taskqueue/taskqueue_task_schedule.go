package taskqueue

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
)

const (
	missedTaskDuration   = 3 * time.Minute
	oldestTaskDuration   = 7 * 24 * time.Hour
	maxTaskCreateForAJob = 5
)

func (q *taskQueue) scanTasksForRunning(ctx context.Context) ([]*entity.Task, error) {
	timeNow := timeutil.NowUTC()
	var allSchedTasks []*entity.Task
	err := transaction.Execute(ctx, q.db, func(db database.Tx) error {
		// Load active cron jobs
		jobSettings, _, err := q.settingRepo.List(ctx, db, nil,
			bunex.SelectWhere("setting.type = ?", base.SettingTypeCronJob),
			bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
			bunex.SelectFor("UPDATE OF setting"),
			bunex.SelectRelation("Tasks",
				bunex.SelectWhereGroup(
					bunex.SelectWhere("task.status = ?", base.TaskStatusNotStarted),
					bunex.SelectWhere("task.run_at >= ?", timeNow.Add(-missedTaskDuration)),
				),
				// Tasks need retry
				bunex.SelectWhereOrGroup(
					bunex.SelectWhere("task.status = ?", base.TaskStatusFailed),
					bunex.SelectWhere("task.max_retry > task.retry"),
					bunex.SelectWhere("task.run_at >= ?", timeNow.Add(-oldestTaskDuration)),
				),
			),
		)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if len(jobSettings) == 0 {
			return nil
		}

		allSchedTasks = make([]*entity.Task, 0, 10) //nolint:mnd
		allNewTasks := make([]*entity.Task, 0, 10)  //nolint:mnd
		for _, setting := range jobSettings {
			schedTasks, newTasks, err := q.getTasksToRun(setting, setting.Tasks, timeNow)
			if err != nil {
				return apperrors.Wrap(err)
			}
			allSchedTasks = append(allSchedTasks, schedTasks...)
			allNewTasks = append(allNewTasks, newTasks...)
		}

		err = q.taskRepo.UpsertMulti(ctx, db, allNewTasks,
			entity.TaskUpsertingConflictCols, entity.TaskUpsertingUpdateCols)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return allSchedTasks, nil
}

//nolint:gocognit
func (q *taskQueue) getTasksToRun(
	jobSetting *entity.Setting,
	lastTasks []*entity.Task,
	timeNow time.Time,
) (schedTasks []*entity.Task, newTasks []*entity.Task, err error) {
	// Find the closest future task to run
	var lastTaskRunAt time.Time
	notStartedTaskCount := 0
	var missedTask, nearestTask *entity.Task
	for _, task := range lastTasks {
		if !task.DeletedAt.IsZero() {
			continue
		}
		if task.Status == base.TaskStatusFailed && task.MaxRetry == task.Retry {
			schedTasks = append(schedTasks, task)
			continue
		}
		if task.Status == base.TaskStatusNotStarted && !task.RunAt.Before(timeNow) {
			if nearestTask == nil {
				nearestTask = task
			}
			notStartedTaskCount++
		}
		if lastTaskRunAt.IsZero() || task.RunAt.After(lastTaskRunAt) {
			lastTaskRunAt = task.RunAt
		}
		if task.RunAt.Before(timeNow) {
			if missedTask == nil || missedTask.RunAt.Before(task.RunAt) {
				missedTask = task
			}
		}
	}

	if !q.shouldRunMissedTask(missedTask, nearestTask, timeNow) {
		missedTask = nil
	}

	if notStartedTaskCount >= maxTaskCreateForAJob/2 {
		schedTasks = append(schedTasks, gofn.ToSliceSkippingNil(missedTask, nearestTask)...)
		return schedTasks, nil, nil
	}

	cronJob, err := jobSetting.AsCronJob()
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}
	cronSched, err := cronJob.ParseCron()
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}

	nextRunAt := gofn.Coalesce(lastTaskRunAt, cronJob.InitialTime)
	farthestRunAt := timeNow.AddDate(1, 0, 0)

	for {
		nextRunAt = cronSched.Next(nextRunAt)
		if nextRunAt.Before(timeNow) {
			continue
		}
		if nextRunAt.After(farthestRunAt) || len(newTasks) >= maxTaskCreateForAJob {
			break
		}

		task := &entity.Task{
			ID:             gofn.Must(ulid.NewStringULID()),
			JobID:          jobSetting.ID,
			Type:           base.TaskType(jobSetting.Kind),
			Status:         base.TaskStatusNotStarted,
			Priority:       cronJob.Priority,
			MaxRetry:       cronJob.MaxRetry,
			RetryDelaySecs: cronJob.RetryDelaySecs,
			Version:        entity.CurrentTaskVersion,
			RunAt:          nextRunAt,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}
		newTasks = append(newTasks, task)

		if nearestTask == nil {
			nearestTask = task
		}
	}

	if !q.shouldRunMissedTask(missedTask, nearestTask, timeNow) {
		missedTask = nil
	}
	schedTasks = append(schedTasks, gofn.ToSliceSkippingNil(missedTask, nearestTask)...)
	return schedTasks, newTasks, nil
}

func (q *taskQueue) shouldRunMissedTask(missedTask, nearestTask *entity.Task, timeNow time.Time) bool {
	if missedTask == nil {
		return false
	}
	if nearestTask == nil {
		return true
	}
	if nearestTask.RunAt.Sub(timeNow) < timeNow.Sub(missedTask.RunAt) {
		return false
	}
	return true
}
