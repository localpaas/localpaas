package appuc

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

type appHttpSettingsData struct {
	HttpSettings *entity.Setting
}

func (uc *AppUC) loadAppDataForUpdateHttpSettings(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	// TODO: add implementation
	return nil
}

//nolint:unparam
func (uc *AppUC) prepareUpdatingAppHttpSettings(
	req *appdto.UpdateAppSettingsReq,
	timeNow time.Time,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) error {
	app := data.App
	setting := data.HttpSettingsData.HttpSettings

	if setting == nil {
		setting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppHttp,
			CreatedAt: timeNow,
		}
	}
	setting.UpdatedAt = timeNow
	setting.Status = base.SettingStatusActive
	setting.ExpireAt = time.Time{}

	httpReq := req.HttpSettings
	setting.MustSetData(&entity.AppHttpSettings{
		Enabled: httpReq.Enabled,
	})

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
	return nil
}

func (uc *AppUC) applyAppHttpSettings(
	_ context.Context,
	_ database.IDB,
	_ *appdto.UpdateAppSettingsReq,
	_ *updateAppSettingsData,
) error {
	// TODO: add implementation
	return nil
}
