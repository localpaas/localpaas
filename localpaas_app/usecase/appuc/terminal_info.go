package appuc

import (
	"context"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *UC) GetTerminalInfo(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.GetTerminalInfoReq,
) (_ *appdto.GetTerminalInfoResp, err error) {
	app, err := uc.appService.LoadApp(ctx, uc.db, req.ProjectID, req.AppID, true, true,
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
		bunex.SelectRelation("Project",
			bunex.SelectExcludeColumns(entity.ProjectDefaultExcludeColumns...),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if app.ServiceID == "" {
		return nil, apperrors.NewUnavailable("App service").
			WithMsgLog("service not exist for app")
	}

	terminalEnabled := true
	featureSetting, err := uc.settingRepo.GetSingle(ctx, uc.db, app.GetSettingScope(),
		base.SettingTypeAppFeatures, true)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.Wrap(err)
	}
	if featureSetting != nil {
		featureSettings := featureSetting.MustAsAppFeatureSettings()
		if featureSettings.TerminalSettings != nil {
			terminalEnabled = featureSettings.TerminalSettings.Enabled
		}
	}
	if !terminalEnabled {
		return &appdto.GetTerminalInfoResp{
			Data: &appdto.TerminalInfoDataResp{Enabled: false},
		}, nil
	}

	return &appdto.GetTerminalInfoResp{
		Data: &appdto.TerminalInfoDataResp{
			Enabled:         true,
			SupportedShells: appdto.SupportedShells,
		},
	}, nil
}
