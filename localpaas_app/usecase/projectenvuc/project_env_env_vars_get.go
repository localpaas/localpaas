package projectenvuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
)

func (uc *ProjectEnvUC) GetProjectEnvEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.GetProjectEnvEnvVarsReq,
) (*projectenvdto.GetProjectEnvEnvVarsResp, error) {
	projectEnv, err := uc.projectEnvRepo.GetByID(ctx, uc.db, req.ProjectEnvID,
		bunex.SelectRelation("EnvVarsSettings"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if projectEnv.ProjectID != req.ProjectID {
		return nil, apperrors.New(apperrors.ErrUnauthorized)
	}

	envVars, err := projectEnv.GetEnvVars()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := projectenvdto.TransformProjectEnvEnvVars(envVars)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.GetProjectEnvEnvVarsResp{
		Data: resp,
	}, nil
}
