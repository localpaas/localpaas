package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) GetAppServiceSpec(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.GetAppServiceSpecReq,
) (*appdto.GetAppServiceSpecResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ProjectID, req.AppID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	service, err := uc.appService.ServiceInspect(ctx, app.ServiceID, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := appdto.TransformAppServiceSpec(service)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.GetAppServiceSpecResp{
		Data: resp,
	}, nil
}
