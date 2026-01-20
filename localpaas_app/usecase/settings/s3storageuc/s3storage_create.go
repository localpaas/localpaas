package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
)

const (
	currentSettingType    = base.SettingTypeS3Storage
	currentSettingVersion = entity.CurrentS3StorageVersion
)

func (uc *S3StorageUC) CreateS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.CreateS3StorageReq,
) (*s3storagedto.CreateS3StorageResp, error) {
	req.Type = currentSettingType
	resp, err := settings.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &settings.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		Version:       currentSettingVersion,
		PrepareCreation: func(
			ctx context.Context,
			db database.Tx,
			data *settings.CreateSettingData,
			pData *settings.PersistingSettingCreationData,
		) error {
			err := pData.Setting.SetData(&entity.S3Storage{
				AccessKeyID: req.AccessKeyID,
				SecretKey:   entity.NewEncryptedField(req.SecretKey),
				Region:      req.Region,
				Bucket:      req.Bucket,
				Endpoint:    req.Endpoint,
			})
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.CreateS3StorageResp{
		Data: resp.Data,
	}, nil
}
