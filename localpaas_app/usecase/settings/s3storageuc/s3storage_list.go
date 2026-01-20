package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
)

func (uc *S3StorageUC) ListS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.ListS3StorageReq,
) (*s3storagedto.ListS3StorageResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := s3storagedto.TransformS3Storages(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.ListS3StorageResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
