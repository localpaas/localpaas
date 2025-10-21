package projectuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

func (uc *ProjectUC) ListProjectBase(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.ListProjectBaseReq,
) (*projectdto.ListProjectBaseResp, error) {
	var listOpts []bunex.SelectQueryOption

	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("project.status IN (?)", bunex.In(req.Status)),
		)
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhere("project.name ILIKE ?", keyword),
		)
	}

	projects, pagingMeta, err := uc.projectRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.ListProjectBaseResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: projectdto.TransformProjectsBase(projects),
	}, nil
}
