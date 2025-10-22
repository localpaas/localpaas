package projectenvuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
)

func (uc *ProjectEnvUC) ListProjectEnv(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.ListProjectEnvReq,
) (*projectenvdto.ListProjectEnvResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("project_env.project_id = ?", req.ProjectID),
	}

	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("project_env.status IN (?)", bunex.In(req.Status)),
		)
	}
	// Filter by search keyword
	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("project_env.name ILIKE ?", keyword),
			),
		)
	}

	projectEnvs, paging, err := uc.projectEnvRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := projectenvdto.TransformProjectEnvs(projectEnvs)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.ListProjectEnvResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}
