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

	ParsedDeploymentSettings *entity.AppDeploymentSettings
	RegistryAuth             *entity.RegistryAuth
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
			deployment, err := deploymentSettings.ParseAppDeploymentSettings()
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to parse app deployment settings")
			}
			deploymentData.ParsedDeploymentSettings = deployment
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
	setting := data.DeploymentData.DeploymentSettings

	if setting == nil {
		setting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppDeployment,
			CreatedAt: timeNow,
		}
	}
	setting.UpdatedAt = timeNow
	setting.ExpireAt = time.Time{}
	setting.Status = base.SettingStatusActive

	deploymentSettings := data.DeploymentData.ParsedDeploymentSettings
	if deploymentSettings == nil {
		deploymentSettings = &entity.AppDeploymentSettings{
			Setting:     setting,
			ImageSource: &entity.DeploymentImageSource{},
			CodeSource:  &entity.DeploymentCodeSource{},
		}
		data.DeploymentData.ParsedDeploymentSettings = deploymentSettings
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
	setting.MustSetData(deploymentSettings)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
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
		Deployment:              deploymentData.ParsedDeploymentSettings,
		ImageSourceRegistryAuth: deploymentData.RegistryAuth,
	})
	if err != nil {
		return apperrors.NewInfra(err)
	}
	return nil
}
