package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) UpdateAppSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppSettingsReq,
) (*appdto.UpdateAppSettingsResp, error) {
	var data *updateAppSettingsData
	var persistingData *persistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data = &updateAppSettingsData{}
		err := uc.loadAppSettingsDataForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		err = uc.preparePersistingAppSettings(req, data, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.applyAppSettings(ctx, db, req, data, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	err = uc.postTransactionAppSettings(ctx, uc.db, req, data, persistingData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppSettingsResp{}, nil
}

type updateAppSettingsData struct {
	App *entity.App

	EnvVarsData      appEnvVarsData
	DeploymentData   appDeploymentData
	HttpSettingsData appHttpSettingsData
	Errors           []string // stores errors
	Warnings         []string // stores warnings
}

func (uc *AppUC) loadAppSettingsDataForUpdate(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	var targetTypes []base.SettingType
	switch {
	case req.EnvVars != nil:
		targetTypes = append(targetTypes, base.SettingTypeEnvVar)
	case req.DeploymentSettings != nil:
		targetTypes = append(targetTypes, base.SettingTypeAppDeployment)
	case req.HttpSettings != nil:
		targetTypes = append(targetTypes, base.SettingTypeAppHttp)
	}

	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type IN (?)", bunex.In(targetTypes)),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app

	updateMatched := true
	for _, setting := range app.Settings {
		switch setting.Type { //nolint:exhaustive
		case base.SettingTypeEnvVar:
			data.EnvVarsData.EnvVars = setting
			updateMatched = req.EnvVars == nil || req.EnvVars.UpdateVer == setting.UpdateVer
		case base.SettingTypeAppDeployment:
			data.DeploymentData.DeploymentSettings = setting
			updateMatched = req.DeploymentSettings == nil || req.DeploymentSettings.UpdateVer == setting.UpdateVer
		case base.SettingTypeAppHttp:
			data.HttpSettingsData.HttpSettings = setting
			updateMatched = req.HttpSettings == nil || req.HttpSettings.UpdateVer == setting.UpdateVer
		}
	}
	if !updateMatched {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}

	switch {
	case req.EnvVars != nil:
		err = uc.loadAppDataForUpdateEnvVars(ctx, db, req, data)
	case req.DeploymentSettings != nil:
		err = uc.loadAppDataForUpdateDeploymentSettings(ctx, db, req, data)
	case req.HttpSettings != nil:
		err = uc.loadAppDataForUpdateHttpSettings(ctx, db, req, data)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *AppUC) preparePersistingAppSettings(
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) (err error) {
	timeNow := timeutil.NowUTC()

	switch {
	case req.EnvVars != nil:
		err = uc.prepareUpdatingAppEnvVars(req, timeNow, data, persistingData)
	case req.DeploymentSettings != nil:
		err = uc.prepareUpdatingAppDeploymentSettings(req, timeNow, data, persistingData)
	case req.HttpSettings != nil:
		err = uc.prepareUpdatingAppHttpSettings(req, timeNow, data, persistingData)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *AppUC) applyAppSettings(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) (err error) {
	switch {
	case req.EnvVars != nil:
		err = uc.applyAppEnvVars(ctx, db, req, data, persistingData)
	case req.DeploymentSettings != nil:
		err = uc.applyAppDeploymentSettings(ctx, db, req, data, persistingData)
	case req.HttpSettings != nil:
		err = uc.applyAppHttpSettings(ctx, db, req, data, persistingData)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *AppUC) postTransactionAppSettings(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) (err error) {
	switch {
	case req.EnvVars != nil:
		err = uc.postTransactionAppEnvVars(ctx, db, req, data, persistingData)
	case req.DeploymentSettings != nil:
		err = uc.postTransactionAppDeploymentSettings(ctx, db, req, data, persistingData)
	case req.HttpSettings != nil:
		err = uc.postTransactionAppHttpSettings(ctx, db, req, data, persistingData)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
