package appuc

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

const (
	// TODO: make this configurable
	defaultDeploymentTimeout = 60 * time.Minute
)

func (uc *AppUC) UpdateAppDeploymentSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppDeploymentSettingsReq,
) (*appdto.UpdateAppDeploymentSettingsResp, error) {
	var data *updateAppDeploymentSettingsData
	var persistingData *persistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data = &updateAppDeploymentSettingsData{}
		err := uc.loadAppDeploymentSettingsForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		err = uc.prepareUpdatingAppDeploymentSettings(req, data, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	err = uc.postTransactionAppDeploymentSettings(ctx, persistingData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppDeploymentSettingsResp{}, nil
}

type updateAppDeploymentSettingsData struct {
	App                    *entity.App
	DeploymentSettings     *entity.Setting
	CurrDeploymentSettings *entity.AppDeploymentSettings
	RegistryAuth           *entity.Setting
	Errors                 []string // stores errors
	Warnings               []string // stores warnings
}

func (uc *AppUC) loadAppDeploymentSettingsForUpdate(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppDeploymentSettingsReq,
	data *updateAppDeploymentSettingsData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppDeployment),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app
	data.DeploymentSettings, _ = gofn.First(app.Settings)

	deploymentSettings := data.DeploymentSettings
	if deploymentSettings != nil && deploymentSettings.UpdateVer != req.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}

	// Parse the current deployment settings
	if deploymentSettings != nil {
		if deploymentSettings.IsActive() && !deploymentSettings.IsExpired() {
			settingData, err := deploymentSettings.AsAppDeploymentSettings()
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to parse app deployment settings")
			}
			data.CurrDeploymentSettings = settingData
		}
	}

	// Loads registry auth if needs to
	imageSource := req.ImageSource
	if imageSource != nil && imageSource.RegistryAuth.ID != "" {
		registryAuth, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeRegistryAuth,
			imageSource.RegistryAuth.ID, true)
		if err != nil {
			return apperrors.Wrap(err)
		}
		data.RegistryAuth = registryAuth
	}

	return nil
}

func (uc *AppUC) prepareUpdatingAppDeploymentSettings(
	req *appdto.UpdateAppDeploymentSettingsReq,
	data *updateAppDeploymentSettingsData,
	persistingData *persistingAppData,
) error {
	app := data.App
	settings := data.DeploymentSettings
	timeNow := timeutil.NowUTC()

	if settings == nil {
		settings = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppDeployment,
			CreatedAt: timeNow,
			Version:   entity.CurrentAppDeploymentSettingsVersion,
		}
		data.DeploymentSettings = settings
	}
	settings.UpdateVer++
	settings.UpdatedAt = timeNow
	settings.ExpireAt = time.Time{}
	settings.Status = base.SettingStatusActive

	newDeploymentSettings, err := uc.buildNewAppDeploymentSettings(req, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	settings.MustSetData(newDeploymentSettings)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, settings)

	// Create a deployment and a task for it
	deployment := &entity.Deployment{
		ID:        gofn.Must(ulid.NewStringULID()),
		AppID:     app.ID,
		Settings:  newDeploymentSettings,
		Status:    base.DeploymentStatusNotStarted,
		Version:   entity.CurrentDeploymentVersion,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	persistingData.UpsertingDeployments = append(persistingData.UpsertingDeployments, deployment)

	deploymentTask := &entity.Task{
		ID:     gofn.Must(ulid.NewStringULID()),
		Type:   base.TaskTypeAppDeploy,
		Status: base.TaskStatusNotStarted,
		Config: entity.TaskConfig{
			Priority: base.TaskPriorityDefault,
			Timeout:  timeutil.Duration(defaultDeploymentTimeout),
		},
		Version:   entity.CurrentTaskVersion,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	err = deploymentTask.SetArgs(&entity.TaskAppDeployArgs{
		Deployment: entity.ObjectID{ID: deployment.ID},
	})
	if err != nil {
		return apperrors.Wrap(err)
	}
	persistingData.UpsertingTasks = append(persistingData.UpsertingTasks, deploymentTask)

	return nil
}

func (uc *AppUC) buildNewAppDeploymentSettings(
	req *appdto.UpdateAppDeploymentSettingsReq,
	data *updateAppDeploymentSettingsData,
) (*entity.AppDeploymentSettings, error) {
	newDeploymentSettings := data.CurrDeploymentSettings
	if newDeploymentSettings == nil {
		newDeploymentSettings = &entity.AppDeploymentSettings{}
	}

	if req.ImageSource != nil {
		err := copier.Copy(&newDeploymentSettings.ImageSource, req.ImageSource)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	if req.RepoSource != nil {
		err := copier.Copy(&newDeploymentSettings.RepoSource, req.RepoSource)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	if req.TarballSource != nil {
		err := copier.Copy(&newDeploymentSettings.TarballSource, req.TarballSource)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	if req.PreDeploymentCommand != nil {
		newDeploymentSettings.PreDeploymentCommand = req.PreDeploymentCommand
	}
	if req.PostDeploymentCommand != nil {
		newDeploymentSettings.PostDeploymentCommand = req.PostDeploymentCommand
	}

	return newDeploymentSettings, nil
}

func (uc *AppUC) postTransactionAppDeploymentSettings(
	ctx context.Context,
	persistingData *persistingAppData,
) error {
	for _, task := range persistingData.UpsertingTasks {
		err := uc.taskQueue.ScheduleTask(ctx, task)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}
	return nil
}
