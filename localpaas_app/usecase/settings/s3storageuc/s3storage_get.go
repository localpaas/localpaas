package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
)

func (uc *S3StorageUC) GetS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.GetS3StorageReq,
) (*s3storagedto.GetS3StorageResp, error) {
	req.Type = currentSettingType
	setting, err := settings.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &settings.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsS3Storage().MustDecrypt()
	resp, err := s3storagedto.TransformS3Storage(setting, req.ObjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.GetS3StorageResp{
		Data: resp,
	}, nil
}
