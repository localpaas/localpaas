package projectuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

func (uc *ProjectUC) GetProject(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.GetProjectReq,
) (*projectdto.GetProjectResp, error) {
	project, err := uc.projectRepo.GetByID(ctx, uc.db, req.ID,
		bunex.SelectRelation("Tags", bunex.SelectOrder("display_order")),
		bunex.SelectRelation("Envs", bunex.SelectOrder("display_order")),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := projectdto.TransformProject(project)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.GetProjectResp{
		Data: resp,
	}, nil
}
