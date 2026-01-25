package appuc

import (
	"context"
	"time"

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
	"github.com/localpaas/localpaas/localpaas_app/service/nginxservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) UpdateAppHttpSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppHttpSettingsReq,
) (*appdto.UpdateAppHttpSettingsResp, error) {
	var data *updateAppHttpSettingsData
	var persistingData *persistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data = &updateAppHttpSettingsData{}
		err := uc.loadAppHttpSettingsForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		err = uc.prepareUpdatingAppHttpSettings(req, data, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.applyAppHttpSettings(ctx, db, data)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppHttpSettingsResp{}, nil
}

type updateAppHttpSettingsData struct {
	App              *entity.App
	HttpSettings     *entity.Setting
	CurrHttpSettings *entity.AppHttpSettings
	Errors           []string // stores errors
	Warnings         []string // stores warnings
}

func (uc *AppUC) loadAppHttpSettingsForUpdate(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppHttpSettingsReq,
	data *updateAppHttpSettingsData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppHttp),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app
	data.HttpSettings, _ = gofn.First(app.Settings)

	if data.HttpSettings != nil && data.HttpSettings.UpdateVer != req.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}

	return nil
}

func (uc *AppUC) prepareUpdatingAppHttpSettings(
	req *appdto.UpdateAppHttpSettingsReq,
	data *updateAppHttpSettingsData,
	persistingData *persistingAppData,
) error {
	app := data.App
	setting := data.HttpSettings
	timeNow := timeutil.NowUTC()

	if setting == nil {
		setting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppHttp,
			CreatedAt: timeNow,
			Version:   entity.CurrentAppHttpSettingsVersion,
		}
		data.HttpSettings = setting
	}
	setting.UpdateVer++
	setting.UpdatedAt = timeNow
	setting.Status = base.SettingStatusActive
	setting.ExpireAt = time.Time{}

	newHttpSettings, err := uc.buildNewAppHttpSettings(req, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	setting.MustSetData(newHttpSettings)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
	return nil
}

//nolint:unparam
func (uc *AppUC) buildNewAppHttpSettings(
	req *appdto.UpdateAppHttpSettingsReq,
	data *updateAppHttpSettingsData,
) (*entity.AppHttpSettings, error) {
	newHttpSettings := data.CurrHttpSettings
	if newHttpSettings == nil {
		newHttpSettings = &entity.AppHttpSettings{}
	}

	newHttpSettings.Enabled = req.Enabled
	newHttpSettings.Domains = gofn.MapSlice(req.Domains, func(r *appdto.DomainReq) *entity.AppDomain {
		return r.ToEntity()
	})

	return newHttpSettings, nil
}

func (uc *AppUC) applyAppHttpSettings(
	ctx context.Context,
	db database.IDB,
	data *updateAppHttpSettingsData,
) error {
	appHttpSettings, err := data.HttpSettings.AsAppHttpSettings()
	if err != nil {
		return apperrors.Wrap(err)
	}

	refSettingMap, err := uc.appService.LoadReferenceSettings(ctx, db, data.App, data.HttpSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	allSslIDs := appHttpSettings.GetInUseSslCertIDs()
	err = uc.appService.EnsureSslConfigFiles(allSslIDs, false, refSettingMap)
	if err != nil {
		return apperrors.Wrap(err)
	}

	allBasicAuthIDs := appHttpSettings.GetInUseBasicAuthIDs()
	err = uc.appService.EnsureBasicAuthConfigFiles(allBasicAuthIDs, false, refSettingMap)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = uc.nginxService.ApplyAppConfig(ctx, data.App, &nginxservice.AppConfigData{
		HttpSettings:  appHttpSettings,
		RefSettingMap: refSettingMap,
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = uc.networkService.UpdateAppGlobalRoutingNetwork(ctx, data.App, data.HttpSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
