package appuc

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

const (
	// TODO: make this configurable
	defaultDeploymentTimeout = 60 * time.Minute
)

type appDeploymentData struct {
	DeploymentSettings *entity.Setting
	RegistryAuth       *entity.Setting
}

func (uc *AppUC) loadAppDataForUpdateDeploymentSettings(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	deploymentData := data.DeploymentData

	// Parse the current deployment settings
	deploymentSettings := deploymentData.DeploymentSettings
	if deploymentSettings != nil {
		if deploymentSettings.IsActive() && !deploymentSettings.IsExpired() {
			_, err := deploymentSettings.AsAppDeploymentSettings()
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to parse app deployment settings")
			}
		}
	}

	// Loads registry auth if needs to
	imageSource := req.DeploymentSettings.ImageSource
	if imageSource != nil && imageSource.RegistryAuth.ID != "" {
		registryAuth, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeRegistryAuth,
			imageSource.RegistryAuth.ID, true)
		if err != nil {
			return apperrors.Wrap(err)
		}
		deploymentData.RegistryAuth = registryAuth
	}

	return nil
}

func (uc *AppUC) prepareUpdatingAppDeploymentSettings(
	req *appdto.UpdateAppSettingsReq,
	timeNow time.Time,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) error {
	app := data.App
	settings := data.DeploymentData.DeploymentSettings

	if settings == nil {
		settings = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppDeployment,
			CreatedAt: timeNow,
			Version:   entity.CurrentAppDeploymentSettingsVersion,
		}
		data.DeploymentData.DeploymentSettings = settings
	}
	settings.UpdateVer++
	settings.UpdatedAt = timeNow
	settings.ExpireAt = time.Time{}
	settings.Status = base.SettingStatusActive

	newDeploymentSettings := &entity.AppDeploymentSettings{
		ImageSource:   &entity.DeploymentImageSource{},
		RepoSource:    &entity.DeploymentRepoSource{},
		TarballSource: &entity.DeploymentTarballSource{},
	}

	if req.DeploymentSettings.ImageSource != nil {
		if err := copier.Copy(newDeploymentSettings.ImageSource, req.DeploymentSettings.ImageSource); err != nil {
			return apperrors.Wrap(err)
		}
	}
	if req.DeploymentSettings.RepoSource != nil {
		if err := copier.Copy(newDeploymentSettings.RepoSource, req.DeploymentSettings.RepoSource); err != nil {
			return apperrors.Wrap(err)
		}
	}
	if req.DeploymentSettings.TarballSource != nil {
		if err := copier.Copy(newDeploymentSettings.TarballSource, req.DeploymentSettings.TarballSource); err != nil {
			return apperrors.Wrap(err)
		}
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
			Priority:  base.TaskPriorityDefault,
			TimeoutMs: timeutil.NewDurationMs(defaultDeploymentTimeout),
		},
		Version:   entity.CurrentTaskVersion,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	err := deploymentTask.SetArgs(&entity.TaskAppDeployArgs{
		Deployment: entity.ObjectID{ID: deployment.ID},
	})
	if err != nil {
		return apperrors.Wrap(err)
	}
	persistingData.UpsertingTasks = append(persistingData.UpsertingTasks, deploymentTask)

	return nil
}

func (uc *AppUC) applyAppDeploymentSettings(
	ctx context.Context,
	_ database.Tx,
	_ *appdto.UpdateAppSettingsReq,
	_ *updateAppSettingsData,
	_ *persistingAppData,
) error {
	return nil
}

func (uc *AppUC) postTransactionAppDeploymentSettings(
	ctx context.Context,
	_ database.IDB,
	_ *appdto.UpdateAppSettingsReq,
	_ *updateAppSettingsData,
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
