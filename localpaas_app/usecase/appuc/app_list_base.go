package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) ListAppBase(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.ListAppBaseReq,
) (*appdto.ListAppBaseResp, error) {
	var listOpts []bunex.SelectQueryOption

	if req.ProjectID != "" {
		listOpts = append(listOpts,
			bunex.SelectWhere("app.project_id = ?", req.ProjectID),
		)
	}
	if req.ProjectEnvID != "" {
		listOpts = append(listOpts,
			bunex.SelectWhere("app.project_env_id = ?", req.ProjectEnvID),
		)
	}
	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("app.status IN (?)", bunex.In(req.Status)),
		)
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhere("app.name ILIKE ?", keyword),
		)
	}

	apps, pagingMeta, err := uc.appRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.ListAppBaseResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: appdto.TransformAppsBase(apps),
	}, nil
}
