package taskappnotification

import (
	"context"
	"errors"
	"fmt"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/emailservice"
	"github.com/localpaas/localpaas/localpaas_app/service/imservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
)

type Executor struct {
	logger         logging.Logger
	redisClient    rediscache.Client
	appRepo        repository.AppRepo
	settingRepo    repository.SettingRepo
	deploymentRepo repository.DeploymentRepo
	taskInfoRepo   cacherepository.TaskInfoRepo
	appService     appservice.AppService
	userService    userservice.UserService
	emailService   emailservice.EmailService
	imService      imservice.IMService
}

func NewExecutor(
	taskQueue taskqueue.TaskQueue,
	logger logging.Logger,
	redisClient rediscache.Client,
	appRepo repository.AppRepo,
	settingRepo repository.SettingRepo,
	deploymentRepo repository.DeploymentRepo,
	taskInfoRepo cacherepository.TaskInfoRepo,
	appService appservice.AppService,
	userService userservice.UserService,
	emailService emailservice.EmailService,
	imService imservice.IMService,
) *Executor {
	p := &Executor{
		logger:         logger,
		redisClient:    redisClient,
		appRepo:        appRepo,
		settingRepo:    settingRepo,
		deploymentRepo: deploymentRepo,
		taskInfoRepo:   taskInfoRepo,
		appService:     appService,
		userService:    userService,
		emailService:   emailService,
		imService:      imService,
	}
	taskQueue.RegisterExecutor(base.TaskTypeAppNotification, p.execute)
	return p
}

type taskData struct {
	*taskqueue.TaskExecData
	NtfnSettings  *entity.AppNotificationSettings
	Project       *entity.Project
	App           *entity.App
	Deployment    *entity.Deployment
	RefSettingMap map[string]*entity.Setting
}

func (e *Executor) execute(
	ctx context.Context,
	db database.Tx,
	task *taskqueue.TaskExecData,
) (err error) {
	data := &taskData{TaskExecData: task}
	err = e.loadNotificationData(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if data.Project == nil || data.App == nil || data.NtfnSettings == nil {
		return nil
	}

	defer func() {
		if err == nil {
			if r := recover(); r != nil {
				err = apperrors.NewPanic(fmt.Sprintf("%v", r))
			}
		}
	}()

	taskTimeout := data.Task.Config.Timeout.ToDuration()
	if taskTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, taskTimeout)
		defer cancel()
	}

	switch { //nolint:gocritic
	case data.Deployment != nil:
		err = e.notifyForDeployment(ctx, db, data)
	}

	return apperrors.Wrap(err)
}

func (e *Executor) loadNotificationData(
	ctx context.Context,
	db database.Tx,
	data *taskData,
) error {
	task := data.Task
	args, err := task.ArgsAsAppNotification()
	if err != nil {
		return apperrors.Wrap(err)
	}

	app, err := e.appRepo.GetByID(ctx, db, "", args.App.ID,
		bunex.SelectWhere("app.status = ?", base.AppStatusActive),
		bunex.SelectRelation("Project"),
		bunex.SelectWhere("project.status = ?", base.ProjectStatusActive),
	)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil
		}
		return apperrors.Wrap(err)
	}
	data.App = app
	data.Project = app.Project

	// Load notification settings
	setting, err := e.settingRepo.GetSingleByAppObject(ctx, db, base.SettingTypeAppNotification, app, true)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil
		}
		return apperrors.Wrap(err)
	}

	ntfnSettings, err := setting.AsAppNotificationSettings()
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.NtfnSettings = ntfnSettings

	if args.Deployment != nil && args.Deployment.ID != "" && ntfnSettings.HasDeploymentNotificationSettings() {
		deployment, err := e.deploymentRepo.GetByID(ctx, db, app.ID, args.Deployment.ID)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.Wrap(err)
		}
		data.Deployment = deployment
	}

	// Load reference settings
	setting.RefIDs = ntfnSettings.GetRefSettingIDs()
	refSettingMap, err := e.appService.LoadReferenceSettings(ctx, db, app, true, setting)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.RefSettingMap = refSettingMap

	return nil
}
