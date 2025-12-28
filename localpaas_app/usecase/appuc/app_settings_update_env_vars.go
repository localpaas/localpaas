package appuc

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

type appEnvVarsData struct {
	EnvVars *entity.Setting
}

func (uc *AppUC) loadAppDataForUpdateEnvVars(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	return nil
}

//nolint:unparam
func (uc *AppUC) prepareUpdatingAppEnvVars(
	req *appdto.UpdateAppSettingsReq,
	timeNow time.Time,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) error {
	app := data.App
	setting := data.EnvVarsData.EnvVars

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
		Data: gofn.MapSlice(req.EnvVars.EnvVars, func(v *appdto.EnvVarReq) *entity.EnvVar {
			return v.ToEntity()
		}),
	})

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
	return nil
}

func (uc *AppUC) applyAppEnvVars(
	ctx context.Context,
	db database.Tx,
	_ *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
	_ *persistingAppData,
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

func (uc *AppUC) postTransactionAppEnvVars(
	_ context.Context,
	_ database.IDB,
	_ *appdto.UpdateAppSettingsReq,
	_ *updateAppSettingsData,
	_ *persistingAppData,
) error {
	return nil
}
