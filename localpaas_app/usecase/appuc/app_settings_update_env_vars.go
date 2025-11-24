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
	DbEnvVarsSettings *entity.Setting
}

func (uc *AppUC) loadAppDataForUpdateEnvVars(
	ctx context.Context,
	db database.IDB,
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
	dbSetting := data.EnvVarsData.DbEnvVarsSettings

	if dbSetting == nil {
		dbSetting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeEnvVar,
			CreatedAt: timeNow,
		}
	}
	dbSetting.UpdatedAt = timeNow
	dbSetting.ExpireAt = time.Time{}
	dbSetting.Status = base.SettingStatusActive

	dbSetting.MustSetData(&entity.EnvVars{Data: gofn.MapSlice(req.EnvVars, func(v *appdto.EnvVarReq) *entity.EnvVar {
		return v.ToEntity()
	})})

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, dbSetting)
	return nil
}

func (uc *AppUC) applyAppEnvVars(
	ctx context.Context,
	db database.IDB,
	_ *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	app := data.App
	envVars, err := uc.envVarService.BuildAppEnv(ctx, db, app, false)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if len(envVars) == 0 { // we will send no ENV vars to the containers of the app
		envVars = []string{}
	}

	service, _, err := uc.dockerManager.ServiceInspect(ctx, app.ServiceID)
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
