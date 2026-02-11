package taskhealthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
)

type Executor struct {
	logger              logging.Logger
	db                  *database.DB
	redisClient         rediscache.Client
	settingRepo         repository.SettingRepo
	appService          appservice.AppService
	settingService      settingservice.SettingService
	userService         userservice.UserService
	notificationService notificationservice.NotificationService
}

func NewExecutor(
	logger logging.Logger,
	db *database.DB,
	taskQueue queue.TaskQueue,
	redisClient rediscache.Client,
	settingRepo repository.SettingRepo,
	appService appservice.AppService,
	settingService settingservice.SettingService,
	userService userservice.UserService,
	notificationService notificationservice.NotificationService,
) *Executor {
	p := &Executor{
		logger:              logger,
		db:                  db,
		redisClient:         redisClient,
		settingRepo:         settingRepo,
		appService:          appService,
		settingService:      settingService,
		userService:         userService,
		notificationService: notificationService,
	}
	taskQueue.RegisterExecutor(base.TaskTypeHealthcheck, p.execute)
	return p
}

type taskData struct {
	*queue.TaskExecData
	HealthcheckSetting *entity.Setting
	Healthcheck        *entity.Healthcheck
	Output             *entity.TaskHealthcheckOutput
	Project            *entity.Project
	App                *entity.App
	NtfnMsgData        *notificationservice.BaseMsgDataHealthcheckNotification
}

func (e *Executor) execute(
	ctx context.Context,
	db database.Tx,
	execData *queue.TaskExecData,
) (err error) {
	task := execData.Task
	data := &taskData{
		TaskExecData: execData,
	}
	data.HealthcheckSetting = data.ObjectMap[task.TargetID].(*entity.Setting) //nolint
	data.Healthcheck = data.HealthcheckSetting.MustAsHealthcheck()
	data.Output = &entity.TaskHealthcheckOutput{}
	data.Project = data.HealthcheckSetting.BelongToProject
	data.App = data.HealthcheckSetting.BelongToApp
	if data.App != nil {
		data.Project = data.App.Project
	}

	defer func() {
		r := recover()
		if err == nil && r != nil {
			err = apperrors.NewPanic(fmt.Sprintf("%v", r))
		}

		task.Status = gofn.If(err == nil, base.TaskStatusDone, base.TaskStatusFailed)
		task.EndedAt = timeutil.NowUTC()
		task.MustSetOutput(data.Output)

		err = e.sendNotification(ctx, db, data)
	}()

	retries := 0
	startTime := time.Now()
	for {
		switch data.Healthcheck.Type { //nolint
		case base.HealthcheckTypeREST:
			err = e.doHealthcheckREST(ctx, data)
		case base.HealthcheckTypeGRPC:
			err = e.doHealthcheckGRPC(ctx, data)
		}
		if err != nil {
			retries++
			if retries > task.Config.MaxRetry {
				break
			}
			task.Config.Retry = retries
			if task.Config.RetryDelay > 0 {
				time.Sleep(task.Config.RetryDelay.ToDuration())
			}
			if time.Since(startTime)+5*time.Second > data.Healthcheck.Interval.ToDuration() {
				break
			}
		} else {
			break
		}
	}

	return nil
}
