package appuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) GetAppEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.GetAppEnvVarsReq,
) (*appdto.GetAppEnvVarsResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ProjectID, req.AppID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// When get ENV vars for an app, also need all ENV vars of the parent app and project
	objectIDs := gofn.ToSliceSkippingZero(app.ID, app.ParentID, app.ProjectID)

	settings, _, err := uc.settingRepo.List(ctx, uc.db, nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeEnvVar),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhere("setting.object_id IN (?)", bunex.In(objectIDs)),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := appdto.TransformEnvVars(app, settings)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.GetAppEnvVarsResp{
		Data: resp,
	}, nil
}
