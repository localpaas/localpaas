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
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

func (uc *AppUC) UpdateAppEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppEnvVarsReq,
) (*appdto.UpdateAppEnvVarsResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		envData := &updateAppEnvVarsData{}
		err := uc.loadAppEnvVarsDataForUpdate(ctx, db, req, envData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingAppData{}
		err = uc.preparePersistingAppEnvVars(req, envData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppEnvVarsResp{}, nil
}

type updateAppEnvVarsData struct {
	App *entity.App
}

func (uc *AppUC) loadAppEnvVarsDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppEnvVarsReq,
	data *updateAppEnvVarsData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("EnvVars"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app

	return nil
}

func (uc *AppUC) preparePersistingAppEnvVars(
	req *appdto.UpdateAppEnvVarsReq,
	data *updateAppEnvVarsData,
	persistingData *persistingAppData,
) error {
	timeNow := timeutil.NowUTC()
	app := data.App
	if app.EnvVars == nil {
		app.EnvVars = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			Type:      base.SettingTypeEnvVar,
			CreatedAt: timeNow,
		}
		app.EnvVarsID = app.EnvVars.ID
	}

	app.EnvVars.UpdatedAt = timeNow
	err := app.EnvVars.SetData(&entity.AppEnvVars{Data: req.EnvVars})
	if err != nil {
		return apperrors.Wrap(err)
	}

	app.UpdatedAt = timeNow
	persistingData.UpsertingApps = append(persistingData.UpsertingApps, app)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, app.EnvVars)

	return nil
}
