package projectenvuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
)

func (uc *ProjectEnvUC) ListProjectEnvBase(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.ListProjectEnvBaseReq,
) (*projectenvdto.ListProjectEnvBaseResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("project_env.project_id = ?", req.ProjectID),
	}

	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("project_env.status IN (?)", bunex.In(req.Status)),
		)
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhere("project_env.name ILIKE ?", keyword),
		)
	}

	projects, pagingMeta, err := uc.projectEnvRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.ListProjectEnvBaseResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: projectenvdto.TransformProjectEnvsBase(projects),
	}, nil
}
