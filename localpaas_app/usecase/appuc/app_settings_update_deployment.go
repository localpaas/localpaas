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
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

type appDeploymentData struct {
	DeploymentSettings *entity.Setting
	RegistryAuth       *entity.Setting
}

func (uc *AppUC) loadAppDataForUpdateDeploymentSettings(
	ctx context.Context,
	db database.IDB,
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
	dbDeploymentSettings := data.DeploymentData.DeploymentSettings

	if dbDeploymentSettings == nil {
		dbDeploymentSettings = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppDeployment,
			CreatedAt: timeNow,
		}
		data.DeploymentData.DeploymentSettings = dbDeploymentSettings
	}
	dbDeploymentSettings.UpdatedAt = timeNow
	dbDeploymentSettings.ExpireAt = time.Time{}
	dbDeploymentSettings.Status = base.SettingStatusActive

	deploymentSettings := &entity.AppDeploymentSettings{
		ImageSource: &entity.DeploymentImageSource{},
		CodeSource:  &entity.DeploymentCodeSource{},
	}

	if req.DeploymentSettings.ImageSource != nil {
		if err := copier.Copy(deploymentSettings.ImageSource, req.DeploymentSettings.ImageSource); err != nil {
			return apperrors.Wrap(err)
		}
	}
	if req.DeploymentSettings.CodeSource != nil {
		if err := copier.Copy(deploymentSettings.CodeSource, req.DeploymentSettings.CodeSource); err != nil {
			return apperrors.Wrap(err)
		}
	}
	dbDeploymentSettings.MustSetData(deploymentSettings)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, dbDeploymentSettings)
	return nil
}

func (uc *AppUC) applyAppDeploymentSettings(
	ctx context.Context,
	_ database.IDB,
	_ *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	deploymentData := data.DeploymentData

	// Update service of the app in docker
	_, err := uc.appService.UpdateAppDeployment(ctx, data.App, &appservice.AppDeploymentReq{
		Deployment:              deploymentData.DeploymentSettings,
		ImageSourceRegistryAuth: deploymentData.RegistryAuth,
	})
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
