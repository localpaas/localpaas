package projectuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

func (uc *ProjectUC) GetProjectEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.GetProjectEnvVarsReq,
) (*projectdto.GetProjectEnvVarsResp, error) {
	project, err := uc.projectRepo.GetByID(ctx, uc.db, req.ProjectID,
		bunex.SelectRelation("EnvVars"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	envVars, err := project.ParseEnvVars()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := projectdto.TransformProjectEnvVars(envVars)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.GetProjectEnvVarsResp{
		Data: resp,
	}, nil
}
