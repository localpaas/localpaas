package appuc

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/swarm"
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

func (uc *AppUC) UpdateAppEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppEnvVarsReq,
) (*appdto.UpdateAppEnvVarsResp, error) {
	var data *updateAppEnvVarsData
	var persistingData *persistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data = &updateAppEnvVarsData{}
		err := uc.loadAppEnvVarsForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		uc.prepareUpdatingAppEnvVars(req, data, persistingData)

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.applyAppEnvVars(ctx, db, data)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppEnvVarsResp{}, nil
}

type updateAppEnvVarsData struct {
	App      *entity.App
	EnvVars  *entity.Setting
	Errors   []string // stores errors
	Warnings []string // stores warnings
}

func (uc *AppUC) loadAppEnvVarsForUpdate(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppEnvVarsReq,
	data *updateAppEnvVarsData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeEnvVar),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app

	if len(app.Settings) > 0 {
		data.EnvVars = app.Settings[0]
	}
	if data.EnvVars != nil && data.EnvVars.UpdateVer != req.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}

	return nil
}

func (uc *AppUC) prepareUpdatingAppEnvVars(
	req *appdto.UpdateAppEnvVarsReq,
	data *updateAppEnvVarsData,
	persistingData *persistingAppData,
) {
	app := data.App
	setting := data.EnvVars
	timeNow := timeutil.NowUTC()

	if setting == nil {
		setting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeEnvVar,
			CreatedAt: timeNow,
			Version:   entity.CurrentEnvVarsVersion,
		}
	}
	setting.UpdateVer++
	setting.UpdatedAt = timeNow
	setting.ExpireAt = time.Time{}
	setting.Status = base.SettingStatusActive

	setting.MustSetData(&entity.EnvVars{
		Data: gofn.MapSlice(req.EnvVars, func(v *appdto.EnvVarReq) *entity.EnvVar {
			return v.ToEntity()
		}),
	})

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *AppUC) applyAppEnvVars(
	ctx context.Context,
	db database.Tx,
	data *updateAppEnvVarsData,
) error {
	app := data.App
	envs, err := uc.envVarService.BuildAppEnv(ctx, db, app, false)
	if err != nil {
		return apperrors.Wrap(err)
	}

	envVars := make([]string, 0, len(envs))
	for _, env := range envs {
		envVars = append(envVars, env.ToString("="))
		if env.Error != "" {
			data.Errors = append(data.Errors, env.Error)
		}
	}

	service, err := uc.dockerManager.ServiceInspect(ctx, app.ServiceID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	if service.Spec.TaskTemplate.ContainerSpec == nil {
		service.Spec.TaskTemplate.ContainerSpec = &swarm.ContainerSpec{}
	}
	service.Spec.TaskTemplate.ContainerSpec.Env = envVars

	_, err = uc.dockerManager.ServiceUpdate(ctx, app.ServiceID, &service.Version, &service.Spec)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
