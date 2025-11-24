package appuc

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

type appHttpSettingsData struct {
	DbHttpSettings *entity.Setting
	HttpSettings   *entity.AppHttpSettings
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

func (uc *AppUC) prepareUpdatingAppHttpSettings(
	req *appdto.UpdateAppSettingsReq,
	timeNow time.Time,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) error { //nolint
	app := data.App
	dbSetting := data.HttpSettingsData.DbHttpSettings

	if dbSetting == nil {
		dbSetting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppHttp,
			CreatedAt: timeNow,
		}
		data.HttpSettingsData.DbHttpSettings = dbSetting
	}
	dbSetting.UpdatedAt = timeNow
	dbSetting.Status = base.SettingStatusActive
	dbSetting.ExpireAt = time.Time{}

	// Existing settings
	existingHttpSettings, err := data.HttpSettingsData.DbHttpSettings.ParseAppHttpSettings()
	if err != nil {
		return apperrors.Wrap(err)
	}

	httpReq := req.HttpSettings
	newHttpSettings := &entity.AppHttpSettings{
		Setting: dbSetting,
		Enabled: httpReq.Enabled,
		Domains: gofn.MapSlice(httpReq.Domains, func(r *appdto.DomainReq) *entity.AppDomain {
			return &entity.AppDomain{
				Domain:  r.Domain,
				SslCert: entity.ObjectID{ID: r.SslCert.ID},
			}
		}),
		DomainRedirect:   httpReq.DomainRedirect,
		ContainerPort:    httpReq.ContainerPort,
		ForceHttps:       httpReq.ForceHttps,
		WebsocketEnabled: httpReq.WebsocketEnabled,
		BasicAuth:        entity.ObjectID{ID: httpReq.BasicAuth.ID},
	}

	// If `nginxSettings` is not sent from FE, use the existing
	if httpReq.NginxSettings == nil {
		newHttpSettings.NginxSettings = existingHttpSettings.NginxSettings
	} else {
		newHttpSettings.NginxSettings = &entity.NginxSettings{
			Enabled: httpReq.NginxSettings.Enabled,
			RootDirectives: gofn.MapSlice(httpReq.NginxSettings.RootDirectives,
				func(r *appdto.NginxDirectiveReq) *entity.NginxDirective {
					return &entity.NginxDirective{
						Invisible: r.Invisible,
						Directive: r.Directive,
					}
				}),
			ServerBlock: &entity.NginxServerBlock{
				Invisible: httpReq.NginxSettings.ServerBlock.Invisible,
				Directives: gofn.MapSlice(httpReq.NginxSettings.ServerBlock.Directives,
					func(r *appdto.NginxDirectiveReq) *entity.NginxDirective {
						return &entity.NginxDirective{
							Invisible: r.Invisible,
							Directive: r.Directive,
						}
					}),
			},
		}
	}
	data.HttpSettingsData.HttpSettings = newHttpSettings

	dbSetting.MustSetData(newHttpSettings)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, dbSetting)
	return nil
}

func (uc *AppUC) applyAppHttpSettings(
	ctx context.Context,
	_ database.IDB,
	_ *appdto.UpdateAppSettingsReq,
	data *updateAppSettingsData,
) error {
	err := uc.nginxService.ApplyAppConfig(ctx, data.App, data.HttpSettingsData.HttpSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = uc.networkService.UpdateAppGlobalRoutingNetwork(ctx, data.App, data.HttpSettingsData.HttpSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
