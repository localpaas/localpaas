package appuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) UpdateAppSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppSettingsReq,
) (*appdto.UpdateAppSettingsResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data := &updateAppSettingsData{}
		err := uc.loadAppSettingsDataForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingAppData{}
		uc.preparePersistingAppSettings(req, data, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppSettingsResp{}, nil
}

type updateAppSettingsData struct {
	App *entity.App
}

func (uc *AppUC) loadAppSettingsDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	var targetTypes []base.SettingType
	switch {
	case req.EnvVars != nil:
		targetTypes = append(targetTypes, base.SettingTypeEnvVar)
	case req.DeploymentSettings != nil:
		targetTypes = append(targetTypes, base.SettingTypeDeployment)
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

	return nil
}

func (uc *AppUC) preparePersistingAppSettings(
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) {
	timeNow := timeutil.NowUTC()
	app := data.App

	if req.EnvVars != nil {
		setting := app.GetSettingByType(base.SettingTypeEnvVar)
		if setting == nil {
			setting = &entity.Setting{
				ID:        gofn.Must(ulid.NewStringULID()),
				ObjectID:  app.ID,
				Type:      base.SettingTypeEnvVar,
				Status:    base.SettingStatusActive,
				CreatedAt: timeNow,
			}
		}
		setting.UpdatedAt = timeNow
		setting.MustSetData(&entity.EnvVars{Data: gofn.MapSlice(req.EnvVars, func(v *appdto.EnvVarReq) *entity.EnvVar {
			return v.ToEntity()
		})})

		persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
	}

	if req.DeploymentSettings != nil {
		setting := app.GetSettingByType(base.SettingTypeDeployment)
		if setting == nil {
			setting = &entity.Setting{
				ID:        gofn.Must(ulid.NewStringULID()),
				ObjectID:  app.ID,
				Type:      base.SettingTypeDeployment,
				Status:    base.SettingStatusActive,
				CreatedAt: timeNow,
			}
		}
		setting.UpdatedAt = timeNow
		setting.MustSetData(req.DeploymentSettings.ToEntity())

		persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
	}
}
