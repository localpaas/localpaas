package taskcronjobexec

import (
	"context"
	"fmt"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
	"github.com/localpaas/localpaas/services/docker"
)

type Executor struct {
	logger        logging.Logger
	redisClient   rediscache.Client
	settingRepo   repository.SettingRepo
	taskRepo      repository.TaskRepo
	taskLogRepo   repository.TaskLogRepo
	appService    appservice.AppService
	dockerManager *docker.Manager
}

func NewExecutor(
	logger logging.Logger,
	taskQueue taskqueue.TaskQueue,
	redisClient rediscache.Client,
	settingRepo repository.SettingRepo,
	taskRepo repository.TaskRepo,
	taskLogRepo repository.TaskLogRepo,
	appService appservice.AppService,
	dockerManager *docker.Manager,
) *Executor {
	p := &Executor{
		logger:        logger,
		redisClient:   redisClient,
		settingRepo:   settingRepo,
		taskRepo:      taskRepo,
		taskLogRepo:   taskLogRepo,
		appService:    appService,
		dockerManager: dockerManager,
	}
	taskQueue.RegisterExecutor(base.TaskTypeCronJobExec, p.execute)
	return p
}

type taskData struct {
	*taskqueue.TaskExecData
	CronJobSetting *entity.Setting
	CronJob        *entity.CronJob
	App            *entity.App
	Logs           []*realtimelog.LogFrame
}

func (e *Executor) execute(
	ctx context.Context,
	db database.Tx,
	task *taskqueue.TaskExecData,
) (err error) {
	data := &taskData{
		TaskExecData:   task,
		CronJobSetting: task.Task.Job,
		CronJob:        task.Task.Job.MustAsCronJob(),
	}
	data.OnPostExec(func() { e.onPostExec(data) })

	err = e.loadCronJobData(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	defer func() {
		if err == nil {
			if r := recover(); r != nil {
				err = apperrors.NewPanic(fmt.Sprintf("%v", r))
			}
		}
		_ = e.saveLogs(ctx, db, data)
	}()

	switch data.CronJob.CronType { //nolint
	case base.CronJobTypeContainerCommand:
		err = e.cronExecContainerCmd(ctx, data)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) loadCronJobData(
	ctx context.Context,
	db database.Tx,
	data *taskData,
) error {
	if data.CronJob.App.ID != "" {
		app, err := e.appService.LoadApp(ctx, db, "", data.CronJob.App.ID, true, true,
			bunex.SelectRelation("Project"),
		)
		if err != nil {
			return apperrors.Wrap(err)
		}
		data.App = app
	}

	return nil
}

func (e *Executor) saveLogs(
	ctx context.Context,
	db database.Tx,
	taskData *taskData,
) error {
	// Insert data in to DB by chunk to avoid exceeding DBMS limit
	for _, chunk := range gofn.Chunk(taskData.Logs, 10000) { //nolint
		taskLogs := make([]*entity.TaskLog, 0, len(chunk))
		for _, logFrame := range chunk {
			taskLogs = append(taskLogs, &entity.TaskLog{
				TaskID:   taskData.Task.ID,
				TargetID: taskData.CronJobSetting.ID,
				Type:     logFrame.Type,
				Data:     logFrame.Data,
				Ts:       logFrame.Ts,
			})
		}
		err := e.taskLogRepo.InsertMulti(ctx, db, taskLogs)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func (e *Executor) onPostExec(
	data *taskData,
) {
	if data.Task.IsDone() || data.Task.IsFailedCompletely() {
		err := e.createNotificationTask(data)
		if err != nil {
			data.Logs = append(data.Logs, realtimelog.NewOutFrame("Failed to schedule notification sending"+
				" with error: "+err.Error(), nil))
		}
	}
}

//nolint:unparam
func (e *Executor) createNotificationTask(
	data *taskData,
) error {
	ntfnSettings := data.CronJob.Notification
	if ntfnSettings == nil {
		return nil
	}

	task := data.Task
	if task.Status != base.TaskStatusDone && task.Status != base.TaskStatusFailed {
		return nil
	}
	if task.Status == base.TaskStatusDone && ntfnSettings.Success == nil {
		return nil
	}
	if task.Status == base.TaskStatusFailed && ntfnSettings.Failure == nil {
		return nil
	}

	timeNow := timeutil.NowUTC()
	ntfnTask := &entity.Task{
		ID:     gofn.Must(ulid.NewStringULID()),
		Type:   base.TaskTypeCronJobNotification,
		Status: base.TaskStatusNotStarted,
		Config: entity.TaskConfig{
			Priority: base.TaskPriorityDefault,
			Timeout:  timeutil.Duration(base.CronJobNotificationTimeoutDefault),
		},
		Version:   entity.CurrentTaskVersion,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	ntfnTask.MustSetArgs(&entity.TaskCronJobNotificationArgs{
		App:     entity.ObjectID{ID: data.App.ID},
		CronJob: entity.ObjectID{ID: data.CronJobSetting.ID},
	})
	data.ScheduleNextTasks(task)
	return nil
}
