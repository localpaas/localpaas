package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) GetAppHttpSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.GetAppHttpSettingsReq,
) (*appdto.GetAppHttpSettingsResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ProjectID, req.AppID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	settings, _, err := uc.settingRepo.List(ctx, uc.db, "", "", nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeAppHttp),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhere("setting.object_id = ?", app.ID),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	input := &appdto.AppHttpSettingsTransformInput{
		App: app,
	}
	if len(settings) > 0 {
		input.HttpSettings = settings[0]
	}

	if input.HttpSettings != nil && input.HttpSettings.MustAsAppHttpSettings() != nil {
		err = uc.loadAppHttpSettingsReferenceData(ctx, uc.db, input)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	resp, err := appdto.TransformHttpSettings(input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.GetAppHttpSettingsResp{
		Data: resp,
	}, nil
}

func (uc *AppUC) loadAppHttpSettingsReferenceData(
	ctx context.Context,
	db database.IDB,
	input *appdto.AppHttpSettingsTransformInput,
) (err error) {
	var settingIDs []string
	for _, domain := range input.HttpSettings.MustAsAppHttpSettings().Domains {
		if domain.SslCert.ID != "" {
			settingIDs = append(settingIDs, domain.SslCert.ID)
		}
		if domain.BasicAuth.ID != "" {
			settingIDs = append(settingIDs, domain.BasicAuth.ID)
		}
	}

	input.ReferenceSettingMap, err = uc.settingRepo.ListByIDsAsMap(ctx, db, settingIDs, true)
	if err != nil {
		return apperrors.Wrap(err)
	}

	input.DefaultNginxSettings, err = uc.nginxService.GetDefaultNginxConfig()
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
