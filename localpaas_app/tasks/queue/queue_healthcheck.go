package queue

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
)

const (
	taskKeyHealthcheck = "task:healthcheck:lock"
	taskLockMaxRetry   = 3
)

func (q *taskQueue) doHealthcheck(
	ctx context.Context,
) error {
	// Make sure only one worker processes this task at a time
	success, releaser, err := q.healthcheckTaskLock(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if !success { // another worker is doing this task
		return nil
	}
	defer releaser()

	executorFunc := q.taskExecutorMap[base.TaskTypeHealthcheck]
	if executorFunc == nil {
		return apperrors.NewUnavailable("Task executor function for healthcheck")
	}

	objectMap := make(map[string]any, 10) //nolint:mnd
	jobSettings, err := q.loadHealthcheckData(ctx, q.db, objectMap)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = transaction.Execute(ctx, q.db, func(db database.Tx) (err error) {
		timeNow := timeutil.NowUTC()
		allTasks := make([]*entity.Task, 0, len(jobSettings))
		execFuncs := make([]func(ctx context.Context) error, 0, len(jobSettings))

		for _, jobSetting := range jobSettings {
			healthcheck := jobSetting.MustAsHealthcheck()
			taskData := &TaskExecData{
				Task: &entity.Task{
					ID:       gofn.Must(ulid.NewStringULID()),
					TargetID: jobSetting.ID,
					Type:     base.TaskTypeHealthcheck,
					Status:   base.TaskStatusNotStarted,
					Config: entity.TaskConfig{
						MaxRetry:   healthcheck.MaxRetry,
						RetryDelay: healthcheck.RetryDelay,
						Timeout:    healthcheck.Timeout,
					},
					Version:   entity.CurrentTaskVersion,
					StartedAt: timeNow,
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
				},
				ObjectMap: objectMap,
			}
			allTasks = append(allTasks, taskData.Task)
			execFuncs = append(execFuncs, func(ctx context.Context) error {
				return executorFunc(ctx, db, taskData) //nolint:wrapcheck
			})
		}

		// Execute all health check tasks concurrently
		_ = gofn.ExecTasksEx(ctx, 20, false, execFuncs...) //nolint:mnd

		// Save tasks in DB
		err = q.taskRepo.UpsertMulti(ctx, db, allTasks,
			entity.TaskUpsertingConflictCols, entity.TaskUpsertingUpdateCols)
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

func (q *taskQueue) healthcheckTaskLock(ctx context.Context) (bool, func(), error) {
	interval := config.Current.Tasks.Healthcheck.Interval
	retries := 0
	wait := time.Duration(0)
	for {
		success, releaser, err := q.taskService.CreateLock(ctx, taskKeyHealthcheck, interval-time.Second)
		if err != nil {
			if retries >= taskLockMaxRetry {
				return false, nil, apperrors.Wrap(err)
			}
			retries++
			wait += time.Second
			time.Sleep(wait)
			continue
		}
		return success, releaser, nil
	}
}

func (q *taskQueue) loadHealthcheckData(
	ctx context.Context,
	db database.IDB,
	objectMap map[string]any,
) ([]*entity.Setting, error) {
	allJobSettings, _, err := q.settingRepo.List(ctx, db, nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeHealthcheck),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectRelation("BelongToProject"),
		bunex.SelectRelation("BelongToApp.Project"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	timeNowSecs := timeutil.NowUTC().Unix()
	refSettingIDs := make([]string, 0, len(allJobSettings))
	validJobSettings := make([]*entity.Setting, 0, len(allJobSettings))
	for _, jobSetting := range allJobSettings {
		project := jobSetting.BelongToProject
		app := jobSetting.BelongToApp
		if app != nil {
			project = app.Project
		}
		if app != nil && app.Status != base.AppStatusActive {
			continue
		}
		if project != nil && project.Status != base.ProjectStatusActive {
			continue
		}

		healthcheck := jobSetting.MustAsHealthcheck()
		if timeNowSecs%int64(healthcheck.Interval.ToDuration().Seconds()) > 5 { //nolint:mnd
			continue
		}
		validJobSettings = append(validJobSettings, jobSetting)
		objectMap[jobSetting.ID] = jobSetting
		refSettingIDs = append(refSettingIDs, healthcheck.GetRefSettingIDs()...)
	}

	refSettings, err := q.settingRepo.ListByIDs(ctx, db, gofn.ToSet(refSettingIDs), true,
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	for _, refSetting := range refSettings {
		objectMap[refSetting.ID] = refSetting
	}

	return validJobSettings, nil
}
