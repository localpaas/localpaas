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
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) UpdateAppDeploymentSource(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppDeploymentSourceReq,
) (*appdto.UpdateAppDeploymentSourceResp, error) {
	var appData *updateAppDeploymentData
	var persistingData *persistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		appData = &updateAppDeploymentData{}
		err := uc.loadAppForUpdateDeployment(ctx, db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		err = uc.prepareUpdatingAppDeployment(req, appData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		// Update service of the app in docker
		_, err = uc.appService.UpdateAppDeployment(ctx, appData.App, &appservice.AppDeploymentReq{
			Deployment:              appData.DeploymentSettings,
			ImageSourceRegistryAuth: appData.RegistryAuth,
		})
		if err != nil {
			return apperrors.NewInfra(err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppDeploymentSourceResp{}, nil
}

type updateAppDeploymentData struct {
	App                *entity.App
	DeploymentSettings *entity.AppDeploymentSettings
	RegistryAuth       *entity.RegistryAuth
}

func (uc *AppUC) loadAppForUpdateDeployment(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppDeploymentSourceReq,
	data *updateAppDeploymentData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE"),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeDeployment),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if app.Status != base.AppStatusActive {
		return apperrors.NewUnavailable("App").
			WithMsgLog("app '%s' is unavailable", app.Name)
	}
	data.App = app

	// Parse the current deployment settings
	currDeploymentSetting := app.GetSettingByType(base.SettingTypeDeployment)
	if currDeploymentSetting != nil {
		deployment, err := currDeploymentSetting.ParseAppDeploymentSettings()
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to parse app deployment settings")
		}
		data.DeploymentSettings = deployment
	}

	// Loads registry auth if needs to
	if req.ImageSource != nil && req.ImageSource.RegistryAuth.ID != "" {
		registryAuthSetting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeRegistryAuth,
			req.ImageSource.RegistryAuth.ID, true)
		if err != nil {
			return apperrors.Wrap(err)
		}
		registryAuth, err := registryAuthSetting.ParseRegistryAuth(true)
		if err != nil {
			return apperrors.Wrap(err)
		}
		data.RegistryAuth = registryAuth
	}

	return nil
}

func (uc *AppUC) prepareUpdatingAppDeployment(
	req *appdto.UpdateAppDeploymentSourceReq,
	data *updateAppDeploymentData,
	persistingData *persistingAppData,
) error {
	timeNow := timeutil.NowUTC()
	app := data.App

	var setting *entity.Setting
	if data.DeploymentSettings != nil {
		setting = data.DeploymentSettings.Setting
	} else {
		setting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeDeployment,
			CreatedAt: timeNow,
		}
	}
	setting.Status = base.SettingStatusActive
	setting.UpdatedAt = timeNow
	setting.ExpireAt = time.Time{}

	deploymentSettings := data.DeploymentSettings
	if deploymentSettings == nil {
		deploymentSettings = &entity.AppDeploymentSettings{
			Setting:     setting,
			ImageSource: &entity.DeploymentImageSource{},
			CodeSource:  &entity.DeploymentCodeSource{},
		}
		data.DeploymentSettings = deploymentSettings
	}
	if req.ImageSource != nil {
		if err := copier.Copy(deploymentSettings.ImageSource, req.ImageSource); err != nil {
			return apperrors.Wrap(err)
		}
	}
	if req.CodeSource != nil {
		if err := copier.Copy(deploymentSettings.CodeSource, req.CodeSource); err != nil {
			return apperrors.Wrap(err)
		}
	}
	setting.MustSetData(deploymentSettings)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
	return nil
}
