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
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

func (uc *AppUC) UpdateAppSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppSettingsReq,
) (*appdto.UpdateAppSettingsResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		settingsData := &updateAppSettingsData{}
		err := uc.loadAppSettingsDataForUpdate(ctx, db, req, settingsData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingAppData{}
		err = uc.preparePersistingAppSettings(req, settingsData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

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
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("Settings"),
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
) error {
	timeNow := timeutil.NowUTC()
	app := data.App
	if app.Settings == nil {
		app.Settings = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			Type:      base.SettingTypeApp,
			CreatedAt: timeNow,
		}
		app.SettingsID = app.Settings.ID
	}

	app.Settings.UpdatedAt = timeNow
	var settingsData *entity.AppSettings

	// Do a copy fields to fields
	err := copier.Copy(&settingsData, req.Settings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = app.Settings.SetData(settingsData)
	if err != nil {
		return apperrors.Wrap(err)
	}

	app.UpdatedAt = timeNow
	persistingData.UpsertingApps = append(persistingData.UpsertingApps, app)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, app.Settings)
	return nil
}
