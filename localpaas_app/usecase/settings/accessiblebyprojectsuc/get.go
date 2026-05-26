package accessiblebyprojectsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/accessiblebyprojectsuc/accessiblebyprojectsdto"
)

func (uc *UC) GetAccessibleByProjects(
	ctx context.Context,
	auth *basedto.Auth,
	req *accessiblebyprojectsdto.GetAccessibleByProjectsReq,
) (*accessiblebyprojectsdto.GetAccessibleByProjectsResp, error) {
	setting, err := uc.SettingRepo.GetByID(ctx, uc.DB, nil, "", req.SettingID, false,
		bunex.SelectRelation("AccessibleByProjects"),
		bunex.SelectRelation("AccessibleByProjects.Project",
			bunex.SelectExcludeColumns(entity.ProjectDefaultExcludeColumns...),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &accessiblebyprojectsdto.GetAccessibleByProjectsResp{
		Data: accessiblebyprojectsdto.TransformAccessibleByProjects(setting),
	}, nil
}
