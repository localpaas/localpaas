package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
)

func (uc *S3StorageUC) GetS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.GetS3StorageReq,
) (*s3storagedto.GetS3StorageResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, base.SettingTypeS3Storage, req.ID, false,
		bunex.SelectRelation("ObjectAccesses",
			bunex.SelectWhere("acl_permission.subject_type IN (?)", bunex.In([]base.SubjectType{
				base.SubjectTypeProject, base.SubjectTypeApp,
			})),
			bunex.SelectRelation("SubjectProject"),
			bunex.SelectRelation("SubjectApp"),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := s3storagedto.TransformS3Storage(setting, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.GetS3StorageResp{
		Data: resp,
	}, nil
}
