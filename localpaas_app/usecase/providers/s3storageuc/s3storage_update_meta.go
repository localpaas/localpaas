package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/s3storageuc/s3storagedto"
)

func (uc *S3StorageUC) UpdateS3StorageMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.UpdateS3StorageMetaReq,
) (*s3storagedto.UpdateS3StorageMetaResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &providers.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.UpdateS3StorageMetaResp{}, nil
}
