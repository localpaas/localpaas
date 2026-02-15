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
	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/cronjobservice"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
	"github.com/localpaas/localpaas/services/docker"
)

type Executor struct {
	logger              logging.Logger
	db                  *database.DB
	redisClient         rediscache.Client
	settingRepo         repository.SettingRepo
	taskRepo            repository.TaskRepo
	taskLogRepo         repository.TaskLogRepo
	cronJobService      cronjobservice.CronJobService
	appService          appservice.AppService
	settingService      settingservice.SettingService
	userService         userservice.UserService
	notificationService notificationservice.NotificationService
	dockerManager       *docker.Manager
}

func NewExecutor(
	logger logging.Logger,
	db *database.DB,
	taskQueue queue.TaskQueue,
	redisClient rediscache.Client,
	settingRepo repository.SettingRepo,
	taskRepo repository.TaskRepo,
	taskLogRepo repository.TaskLogRepo,
	cronJobService cronjobservice.CronJobService,
	appService appservice.AppService,
	settingService settingservice.SettingService,
	userService userservice.UserService,
	notificationService notificationservice.NotificationService,
	dockerManager *docker.Manager,
) *Executor {
	e := &Executor{
		logger:              logger,
		db:                  db,
		redisClient:         redisClient,
		settingRepo:         settingRepo,
		taskRepo:            taskRepo,
		taskLogRepo:         taskLogRepo,
		cronJobService:      cronJobService,
		appService:          appService,
		settingService:      settingService,
		userService:         userService,
		notificationService: notificationService,
		dockerManager:       dockerManager,
	}
	taskQueue.RegisterExecutor(base.TaskTypeCronJobExec, e.execute)
	return e
}

type taskData struct {
	*queue.TaskExecData
	CronJobSetting *entity.Setting
	CronJob        *entity.CronJob
	Project        *entity.Project
	App            *entity.App
	LogStore       *applog.Store
	RefSettingMap  map[string]*entity.Setting
	NotifMsgData   *notificationservice.BaseMsgDataCronTaskNotification
}

func (e *Executor) execute(
	ctx context.Context,
	db database.Tx,
	task *queue.TaskExecData,
) (err error) {
	data := &taskData{
		TaskExecData:   task,
		CronJobSetting: task.Task.TargetJob,
		CronJob:        task.Task.TargetJob.MustAsCronJob(),
	}
	data.OnPostTransaction(func() { e.onPostTransaction(data) }) //nolint

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
		_ = e.saveLogs(ctx, db, data, true)
	}()

	switch data.CronJob.CronType { //nolint
	case base.CronJobTypeContainerCommand:
		err = e.cronExecContainerCmd(ctx, db, data)
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
	logStoreKey := fmt.Sprintf("cron:%s:exec", data.CronJobSetting.ID)
	data.LogStore = applog.NewLocalStore(logStoreKey)

	if data.CronJob.App.ID != "" {
		app, err := e.appService.LoadApp(ctx, db, "", data.CronJob.App.ID, true, true,
			bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
			bunex.SelectRelation("Project",
				bunex.SelectExcludeColumns(entity.ProjectDefaultExcludeColumns...),
			),
		)
		if err != nil {
			return apperrors.Wrap(err)
		}
		data.App = app
		data.Project = app.Project
	}

	// Load reference settings
	refSettingMap, err := e.settingService.LoadReferenceSettingsFor(ctx, db, data.Project, data.App, true,
		data.CronJobSetting)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.RefSettingMap = refSettingMap

	return nil
}

func (e *Executor) saveLogs(
	ctx context.Context,
	db database.IDB,
	data *taskData,
	addDurationInfo bool,
) error {
	logStore := data.LogStore
	if logStore == nil {
		return nil
	}

	if addDurationInfo {
		duration := timeutil.NowUTC().Sub(data.Task.StartedAt)
		_ = logStore.Add(ctx, applog.NewOutFrame("Cron execution finished in "+duration.String(),
			applog.TsNow))
	}

	logFrames, err := logStore.GetData(ctx, 0)
	if err != nil {
		return apperrors.Wrap(err)
	}
	_ = logStore.Close() //nolint

	// Insert data in to DB by chunk to avoid exceeding DBMS limit
	for _, chunk := range gofn.Chunk(logFrames, 10000) { //nolint
		taskLogs := make([]*entity.TaskLog, 0, len(chunk))
		for _, logFrame := range chunk {
			taskLogs = append(taskLogs, &entity.TaskLog{
				TaskID:   data.Task.ID,
				TargetID: data.CronJobSetting.ID,
				Type:     logFrame.Type,
				Data:     logFrame.Data,
				Ts:       logFrame.Ts,
			})
		}
		err = e.taskLogRepo.InsertMulti(ctx, db, taskLogs)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func (e *Executor) onPostTransaction(
	data *taskData,
) {
	ctx := context.Background()
	db := e.db

	// NOTE: We are now outside the transaction, need to reset some data before using them again
	data.LogStore = applog.NewLocalStore(data.LogStore.Key)

	defer func() {
		_ = e.saveLogs(ctx, db, data, false)
	}()

	if data.Task.IsDone() || data.Task.IsFailedCompletely() {
		err := e.sendNotification(ctx, db, data)
		if err != nil {
			_ = data.LogStore.Add(ctx, applog.NewOutFrame("Failed to send result notification"+
				" with error: "+err.Error(), applog.TsNow))
		}
	}
}
