package projectenvuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
)

func (uc *ProjectEnvUC) GetProjectEnv(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.GetProjectEnvReq,
) (*projectenvdto.GetProjectEnvResp, error) {
	projectEnv, err := uc.projectEnvRepo.GetByID(ctx, uc.db, req.ProjectEnvID,
		bunex.SelectRelation("Apps", bunex.SelectOrder("name")),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if projectEnv.ProjectID != req.ProjectID {
		return nil, apperrors.New(apperrors.ErrUnauthorized)
	}

	resp, err := projectenvdto.TransformProjectEnv(projectEnv)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.GetProjectEnvResp{
		Data: resp,
	}, nil
}
