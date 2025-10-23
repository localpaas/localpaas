package projectuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

func (uc *ProjectUC) GetProjectSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.GetProjectSettingsReq,
) (*projectdto.GetProjectSettingsResp, error) {
	project, err := uc.projectRepo.GetByID(ctx, uc.db, req.ProjectID,
		bunex.SelectRelation("Settings"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	settings, err := project.ParseSettings()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := projectdto.TransformProjectSettings(settings)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.GetProjectSettingsResp{
		Data: resp,
	}, nil
}
