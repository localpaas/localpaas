package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) GetApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.GetAppReq,
) (*appdto.GetAppResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ID,
		bunex.SelectRelation("Tags", bunex.SelectOrder("display_order")),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := appdto.TransformApp(app)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.GetAppResp{
		Data: resp,
	}, nil
}
