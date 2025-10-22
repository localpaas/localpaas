package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) GetAppSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.GetAppSettingsReq,
) (*appdto.GetAppSettingsResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ProjectID, req.AppID,
		bunex.SelectRelation("MainSettings"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	settings, err := app.GetMainSettings()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := appdto.TransformAppSettings(settings)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.GetAppSettingsResp{
		Data: resp,
	}, nil
}
