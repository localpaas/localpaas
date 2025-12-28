package appuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) GetAppSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.GetAppSettingsReq,
) (*appdto.GetAppSettingsResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ProjectID, req.AppID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	objectIDs := []string{app.ID}
	// When get ENV vars for an app, also need all ENV vars of the parent app and project
	if gofn.Contain(req.Type, base.SettingTypeEnvVar) {
		objectIDs = append(objectIDs, gofn.ToSliceSkippingZero(app.ParentID, app.ProjectID)...)
	}

	settings, _, err := uc.settingRepo.List(ctx, uc.db, nil,
		bunex.SelectWhere("setting.type IN (?)", bunex.In(req.Type)),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhere("setting.object_id IN (?)", bunex.In(objectIDs)),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	app.Settings = settings

	input := &appdto.AppSettingsTransformationInput{
		App: app,
	}
	err = uc.loadAppSettingsReferenceData(ctx, uc.db, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := appdto.TransformAppSettings(input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.GetAppSettingsResp{
		Data: resp,
	}, nil
}

func (uc *AppUC) loadAppSettingsReferenceData(
	ctx context.Context,
	db database.IDB,
	input *appdto.AppSettingsTransformationInput,
) (err error) {
	app := input.App

	for _, setting := range app.Settings {
		switch setting.Type { //nolint:exhaustive
		case base.SettingTypeEnvVar:
			input.EnvVars = append(input.EnvVars, setting)

		case base.SettingTypeAppDeployment:
			input.DeploymentSettings = setting

		case base.SettingTypeAppHttp:
			input.HttpSettings = setting
		}
	}

	// Reference data for Http settings
	if input.HttpSettings != nil && input.HttpSettings.MustAsAppHttpSettings() != nil {
		err = uc.loadAppHttpSettingsReferenceData(ctx, db, input)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	// Reference data for deployment settings
	if input.DeploymentSettings != nil {
		err = uc.loadAppDeploymentSettingsReferenceData(ctx, db, input)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func (uc *AppUC) loadAppHttpSettingsReferenceData(
	ctx context.Context,
	db database.IDB,
	input *appdto.AppSettingsTransformationInput,
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

func (uc *AppUC) loadAppDeploymentSettingsReferenceData(
	_ context.Context,
	_ database.IDB,
	_ *appdto.AppSettingsTransformationInput,
) (err error) {
	// TODO: add implementation
	return nil
}
