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

func (uc *AppUC) prepareUpdatingAppHttpSettings(
	req *appdto.UpdateAppSettingsReq,
	timeNow time.Time,
	data *updateAppSettingsData,
	persistingData *persistingAppData,
) error { //nolint
	app := data.App
	dbHttpSettings := data.HttpSettingsData.HttpSettings

	if dbHttpSettings == nil {
		dbHttpSettings = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppHttp,
			CreatedAt: timeNow,
		}
		data.HttpSettingsData.HttpSettings = dbHttpSettings
	}
	dbHttpSettings.UpdatedAt = timeNow
	dbHttpSettings.Status = base.SettingStatusActive
	dbHttpSettings.ExpireAt = time.Time{}

	httpReq := req.HttpSettings
	newHttpSettings := &entity.AppHttpSettings{
		Enabled: httpReq.Enabled,
		Domains: gofn.MapSlice(httpReq.Domains, func(r *appdto.DomainReq) *entity.AppDomain {
			return &entity.AppDomain{
				Enabled:          r.Enabled,
				Domain:           r.Domain,
				DomainRedirect:   r.DomainRedirect,
				SslCert:          entity.ObjectID{ID: r.SslCert.ID},
				ContainerPort:    r.ContainerPort,
				ForceHttps:       r.ForceHttps,
				WebsocketEnabled: r.WebsocketEnabled,
				BasicAuth:        entity.ObjectID{ID: r.BasicAuth.ID},
				NginxSettings: &entity.NginxSettings{
					RootDirectives: gofn.MapSlice(r.NginxSettings.RootDirectives,
						func(r *appdto.NginxDirectiveReq) *entity.NginxDirective {
							return &entity.NginxDirective{
								Hide:      r.Hide,
								Directive: r.Directive,
							}
						}),
					ServerBlock: &entity.NginxServerBlock{
						Hide: r.NginxSettings.ServerBlock.Hide,
						Directives: gofn.MapSlice(r.NginxSettings.ServerBlock.Directives,
							func(r *appdto.NginxDirectiveReq) *entity.NginxDirective {
								return &entity.NginxDirective{
									Hide:      r.Hide,
									Directive: r.Directive,
								}
							}),
					},
				},
			}
		}),
	}

	dbHttpSettings.MustSetData(newHttpSettings)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, dbHttpSettings)
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
