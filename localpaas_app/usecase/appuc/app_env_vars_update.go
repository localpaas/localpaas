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
		err = uc.preparePersistingAppEnvVars(auth, req, envData, persistingData)
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
	App              *entity.App
	ExistingSettings *entity.Setting
}

func (uc *AppUC) loadAppEnvVarsDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppEnvVarsReq,
	data *updateAppEnvVarsData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("EnvVarsSettings"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app

	if len(app.EnvVarsSettings) > 0 {
		data.ExistingSettings = app.EnvVarsSettings[0]
	}

	return nil
}

func (uc *AppUC) preparePersistingAppEnvVars(
	auth *basedto.Auth,
	req *appdto.UpdateAppEnvVarsReq,
	data *updateAppEnvVarsData,
	persistingData *persistingAppData,
) error {
	timeNow := timeutil.NowUTC()
	app := data.App
	settings := data.ExistingSettings
	if settings == nil {
		settings = &entity.Setting{
			ID:         gofn.Must(ulid.NewStringULID()),
			TargetType: base.SettingTargetEnvVar,
			TargetID:   app.ID,
			CreatedAt:  timeNow,
			CreatedBy:  auth.User.ID,
		}
	}

	settings.UpdatedAt = timeNow
	settings.UpdatedBy = auth.User.ID

	err := settings.SetData(&entity.AppEnvVars{Data: req.EnvVars})
	if err != nil {
		return apperrors.Wrap(err)
	}

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, settings)
	return nil
}
