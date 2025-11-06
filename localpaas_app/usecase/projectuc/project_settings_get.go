package projectuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

func (uc *ProjectUC) GetProjectSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.GetProjectSettingsReq,
) (*projectdto.GetProjectSettingsResp, error) {
	project, err := uc.projectRepo.GetByID(ctx, uc.db, req.ProjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	objectIDs := []string{project.ID}
	settings, _, err := uc.settingRepo.List(ctx, uc.db, nil,
		bunex.SelectWhere("setting.type IN (?)", bunex.In(req.Type)),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhere("setting.object_id IN (?)", bunex.In(objectIDs)),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	project.Settings = settings

	resp, err := projectdto.TransformProjectSettings(project)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.GetProjectSettingsResp{
		Data: resp,
	}, nil
}
