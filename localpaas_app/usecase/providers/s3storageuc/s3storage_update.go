package s3storageuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/s3storageuc/s3storagedto"
)

func (uc *S3StorageUC) UpdateS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.UpdateS3StorageReq,
) (*s3storagedto.UpdateS3StorageResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &providers.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: gofn.PtrValueOrEmpty(req.Name),
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *providers.UpdateSettingData,
			pData *providers.PersistingSettingData,
		) error {
			setting := pData.Setting
			if req.Name != nil {
				setting.Name = *req.Name
			}
			s3Storage, err := setting.AsS3Storage()
			if err != nil {
				return apperrors.Wrap(err)
			}
			if s3Storage == nil {
				s3Storage = &entity.S3Storage{}
			}
			if req.AccessKeyID != nil {
				s3Storage.AccessKeyID = *req.AccessKeyID
			}
			if req.SecretKey != nil {
				s3Storage.SecretKey = entity.NewEncryptedField(*req.SecretKey)
			}
			if req.Region != nil {
				s3Storage.Region = *req.Region
			}
			if req.Bucket != nil {
				s3Storage.Bucket = *req.Bucket
			}
			if req.Endpoint != nil {
				s3Storage.Endpoint = *req.Endpoint
			}
			if err = setting.SetData(s3Storage); err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.UpdateS3StorageResp{}, nil
}
