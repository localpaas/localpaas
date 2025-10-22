package projectenvuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
)

func (uc *ProjectEnvUC) GetProjectEnvSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.GetProjectEnvSettingsReq,
) (*projectenvdto.GetProjectEnvSettingsResp, error) {
	projectEnv, err := uc.projectEnvRepo.GetByID(ctx, uc.db, req.ProjectEnvID,
		bunex.SelectRelation("MainSettings"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if projectEnv.ProjectID != req.ProjectID {
		return nil, apperrors.New(apperrors.ErrUnauthorized)
	}

	settings, err := projectEnv.GetMainSettings()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := projectenvdto.TransformProjectEnvSettings(settings)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.GetProjectEnvSettingsResp{
		Data: resp,
	}, nil
}
