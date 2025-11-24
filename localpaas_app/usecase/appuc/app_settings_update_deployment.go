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
	DbDeploymentSettings *entity.Setting
	DeploymentSettings   *entity.AppDeploymentSettings
	RegistryAuth         *entity.RegistryAuth
}

func (uc *AppUC) loadAppDataForUpdateDeploymentSettings(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	deploymentData := data.DeploymentData

	// Parse the current deployment settings
	dbDeploymentSettings := deploymentData.DbDeploymentSettings
	if dbDeploymentSettings != nil {
		if dbDeploymentSettings.IsActive() && !dbDeploymentSettings.IsExpired() {
			deployment, err := dbDeploymentSettings.ParseAppDeploymentSettings()
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to parse app deployment settings")
			}
			deploymentData.DeploymentSettings = deployment
		}
	}

	// Loads registry auth if needs to
	imageSource := req.DeploymentSettings.ImageSource
	if imageSource != nil && imageSource.RegistryAuth.ID != "" {
		registryAuthSetting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeRegistryAuth,
			imageSource.RegistryAuth.ID, true)
		if err != nil {
			return apperrors.Wrap(err)
		}
		registryAuth, err := registryAuthSetting.ParseRegistryAuth(true)
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
	dbDeploymentSettings := data.DeploymentData.DbDeploymentSettings

	if dbDeploymentSettings == nil {
		dbDeploymentSettings = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppDeployment,
			CreatedAt: timeNow,
		}
	}
	dbDeploymentSettings.UpdatedAt = timeNow
	dbDeploymentSettings.ExpireAt = time.Time{}
	dbDeploymentSettings.Status = base.SettingStatusActive

	deploymentSettings := data.DeploymentData.DeploymentSettings
	if deploymentSettings == nil {
		deploymentSettings = &entity.AppDeploymentSettings{
			Setting:     dbDeploymentSettings,
			ImageSource: &entity.DeploymentImageSource{},
			CodeSource:  &entity.DeploymentCodeSource{},
		}
		data.DeploymentData.DeploymentSettings = deploymentSettings
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
		return apperrors.NewInfra(err)
	}
	return nil
}
