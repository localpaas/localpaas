package taskqueue

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/gocronqueue"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
)

type TaskExecData struct {
	Task         *entity.Task
	Uncancelable bool
	Canceled     bool
	Done         bool
	onCommand    func(base.TaskCommand, ...any)
}

func (t *TaskExecData) IsCanceled() bool {
	return t.Canceled
}

func (t *TaskExecData) IsDone() bool {
	return t.Done
}

func (t *TaskExecData) OnCommand(fn func(base.TaskCommand, ...any)) {
	// NOTE: do we need to use mutex?
	t.onCommand = fn
}

type TaskExecFunc func(context.Context, database.Tx, *TaskExecData) error

func (q *taskQueue) RegisterExecutor(typ base.TaskType, execFunc TaskExecFunc) {
	if !q.isWorkerMode() {
		return
	}
	if q.taskExecutorMap == nil {
		q.taskExecutorMap = make(map[base.TaskType]gocronqueue.TaskExecFunc, 10) //nolint:mnd
	}
	q.taskExecutorMap[typ] = func(taskID string, payload string) (time.Time, error) {
		return q.executeTask(context.Background(), taskID, payload, execFunc)
	}
}

func (q *taskQueue) ScheduleTask(
	ctx context.Context,
	task *entity.Task,
) error {
	if task.Status == base.TaskStatusDone || task.Status == base.TaskStatusCanceled {
		return nil
	}
	if q.server != nil {
		if err := q.server.ScheduleTask(task, task.ShouldRunAt()); err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	}
	if q.client != nil {
		if err := q.client.ScheduleTask(task, task.ShouldRunAt()); err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	}

	return apperrors.New(apperrors.ErrUnavailable).WithMsgLog("task queue is not initialized")
}

func (q *taskQueue) UnscheduleTask(
	ctx context.Context,
	task *entity.Task,
) error {
	if q.server != nil {
		if err := q.server.UnscheduleTask(task); err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	}
	if q.client != nil {
		if err := q.client.UnscheduleTask(task); err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	}

	return apperrors.New(apperrors.ErrUnavailable).WithMsgLog("task queue is not initialized")
}

func (q *taskQueue) ScheduleTasksForCronJob(
	ctx context.Context,
	db database.Tx,
	jobSetting *entity.Setting,
	unscheduleCurrentTasks bool,
) error {
	if unscheduleCurrentTasks {
		unschedulingTasks, err := q.loadCurrentTasksForUnscheduling(ctx, db, jobSetting)
		if err != nil {
			return apperrors.Wrap(err)
		}
		err = q.taskRepo.UpsertMulti(ctx, db, unschedulingTasks,
			entity.TaskUpsertingConflictCols, entity.TaskUpsertingUpdateCols)
		if err != nil {
			return apperrors.Wrap(err)
		}
		for _, task := range unschedulingTasks {
			if err := q.UnscheduleTask(ctx, task); err != nil {
				return apperrors.Wrap(err)
			}
		}
	}

	if jobSetting.DeletedAt.IsZero() && jobSetting.IsActive() {
		tasks, err := q.createTasks(ctx, db, []string{jobSetting.ID}, q.config.TaskQueue.TaskCreateInterval)
		if err != nil {
			return apperrors.Wrap(err)
		}
		for _, task := range tasks {
			if err := q.ScheduleTask(ctx, task); err != nil {
				return apperrors.Wrap(err)
			}
		}
	}

	return nil
}

func (q *taskQueue) createTasks(ctx context.Context, db database.Tx, jobIDs []string,
	withinDuration time.Duration) ([]*entity.Task, error) {
	opts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", base.SettingTypeCronJob),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectFor("UPDATE OF setting"),
	}
	if len(jobIDs) > 0 {
		opts = append(opts, bunex.SelectWhereIn("setting.id IN (?)", jobIDs...))
	}

	jobSettings, _, err := q.settingRepo.List(ctx, db, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(jobSettings) == 0 {
		return nil, nil
	}

	mapLastRunAt, err := q.queryLastTaskRunAt(ctx, db, entityutil.ExtractIDs(jobSettings))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	timeNow := timeutil.NowUTC()
	allNewTasks := make([]*entity.Task, 0, 20) //nolint:mnd
	for _, jobSetting := range jobSettings {
		cronJob, err := jobSetting.AsCronJob()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		cronSched, err := cronJob.ParseCron()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}

		nextRunAt := gofn.Coalesce(mapLastRunAt[jobSetting.ID], cronJob.InitialTime)
		farthestRunAt := timeNow.Add(withinDuration)

		for {
			nextRunAt = cronSched.Next(nextRunAt)
			if nextRunAt.Before(timeNow) {
				continue
			}
			if nextRunAt.After(farthestRunAt) {
				break
			}

			task := &entity.Task{
				ID:     gofn.Must(ulid.NewStringULID()),
				JobID:  jobSetting.ID,
				Type:   base.TaskType(jobSetting.Kind),
				Status: base.TaskStatusNotStarted,
				Config: entity.TaskConfig{
					Priority:       cronJob.Priority,
					MaxRetry:       cronJob.MaxRetry,
					RetryDelaySecs: cronJob.RetryDelaySecs,
					TimeoutSecs:    cronJob.TimeoutSecs,
				},
				Version:   entity.CurrentTaskVersion,
				RunAt:     nextRunAt,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			}
			allNewTasks = append(allNewTasks, task)
		}
	}

	err = q.taskRepo.UpsertMulti(ctx, db, allNewTasks,
		entity.TaskUpsertingConflictCols, entity.TaskUpsertingUpdateCols)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return allNewTasks, nil
}

func (q *taskQueue) queryLastTaskRunAt(ctx context.Context, db database.IDB, jobIDs []string) (
	map[string]time.Time, error) {
	tasks, _, err := q.taskRepo.List(ctx, db, "", nil,
		bunex.SelectDistinctOn("job_id", "run_at"),
		bunex.SelectWhereIn("job_id IN (?)", jobIDs...),
		bunex.SelectOrder("job_id", "run_at DESC"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	res := make(map[string]time.Time, len(tasks))
	for _, task := range tasks {
		res[task.JobID] = task.RunAt
	}
	return res, nil
}

func (q *taskQueue) loadCurrentTasksForUnscheduling(
	ctx context.Context,
	db database.IDB,
	job *entity.Setting,
) ([]*entity.Task, error) {
	timeNow := timeutil.NowUTC()
	tasks, _, err := q.taskRepo.List(ctx, db, job.ID, nil,
		bunex.SelectFor("UPDATE OF task SKIP LOCKED"),
		bunex.SelectWhere("task.status != ?", base.TaskStatusDone),
		bunex.SelectWhere("task.run_at > ?", timeNow.Add(-10*24*time.Hour)), //nolint scan from 10 days ago
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	unschedulingTasks := make([]*entity.Task, 0, len(tasks))
	for _, task := range tasks {
		if task.CanCancel() {
			task.Status = base.TaskStatusCanceled
			unschedulingTasks = append(unschedulingTasks, task)
			continue
		}
	}

	return unschedulingTasks, nil
}
