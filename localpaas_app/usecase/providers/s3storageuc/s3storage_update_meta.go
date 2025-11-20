package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/s3storageuc/s3storagedto"
)

func (uc *S3StorageUC) UpdateS3StorageMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.UpdateS3StorageMetaReq,
) (*s3storagedto.UpdateS3StorageMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		s3Data := &updateS3StorageData{}
		err := uc.loadS3StorageDataForUpdateMeta(ctx, db, req, s3Data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingS3StorageMeta(req, s3Data)
		return uc.persistS3StorageMeta(ctx, db, s3Data)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.UpdateS3StorageMetaResp{}, nil
}

func (uc *S3StorageUC) loadS3StorageDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *s3storagedto.UpdateS3StorageMetaReq,
	data *updateS3StorageData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeS3Storage, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *S3StorageUC) prepareUpdatingS3StorageMeta(
	req *s3storagedto.UpdateS3StorageMetaReq,
	data *updateS3StorageData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting

	if req.Status != nil {
		setting.Status = *req.Status
	}
	if req.ExpireAt != nil {
		setting.ExpireAt = *req.ExpireAt
	}

	setting.UpdatedAt = timeNow
}

func (uc *S3StorageUC) persistS3StorageMeta(
	ctx context.Context,
	db database.IDB,
	data *updateS3StorageData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
