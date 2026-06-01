package appactionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appactionuc/appactiondto"
)

func (uc *UC) RestartApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *appactiondto.RestartAppReq,
) (*appactiondto.RestartAppResp, error) {
	app, err := uc.appService.LoadApp(ctx, uc.db, req.ProjectID, req.AppID, true, true,
		bunex.SelectRelation("Project",
			bunex.SelectExcludeColumns(entity.ProjectDefaultExcludeColumns...),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	err = uc.dockerManager.ServiceForceUpdate(ctx, app.ServiceID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appactiondto.RestartAppResp{}, nil
}
