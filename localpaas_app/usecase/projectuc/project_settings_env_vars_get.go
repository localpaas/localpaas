package projectuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

func (uc *ProjectUC) GetProjectEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.GetProjectEnvVarsReq,
) (*projectdto.GetProjectEnvVarsResp, error) {
	project, err := uc.projectRepo.GetByID(ctx, uc.db, req.ProjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	settings, _, err := uc.settingRepo.List(ctx, uc.db, nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeEnvVar),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhere("setting.object_id = ?", project.ID),
		bunex.SelectLimit(1),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	var setting *entity.Setting
	if len(settings) > 0 {
		setting = settings[0]
	}

	resp, err := projectdto.TransformEnvVars(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.GetProjectEnvVarsResp{
		Data: resp,
	}, nil
}
