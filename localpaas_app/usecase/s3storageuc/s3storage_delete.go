package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/s3storageuc/s3storagedto"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

func (uc *S3StorageUC) DeleteS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.DeleteS3StorageReq,
) (*s3storagedto.DeleteS3StorageResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		s3storageData := &deleteS3StorageData{}
		err := uc.loadS3StorageDataForDelete(ctx, db, req, s3storageData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingS3StorageData{}
		uc.prepareDeletingS3Storage(s3storageData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.DeleteS3StorageResp{}, nil
}

type deleteS3StorageData struct {
	S3Storage *entity.S3Storage
}

func (uc *S3StorageUC) loadS3StorageDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *s3storagedto.DeleteS3StorageReq,
	data *deleteS3StorageData,
) error {
	s3Storage, err := uc.s3StorageRepo.GetByID(ctx, db, req.ID,
		bunex.SelectFor("UPDATE OF s3_storage"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.S3Storage = s3Storage

	return nil
}

func (uc *S3StorageUC) prepareDeletingS3Storage(
	data *deleteS3StorageData,
	persistingData *persistingS3StorageData,
) {
	s3Storage := data.S3Storage
	s3Storage.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingS3Storages = append(persistingData.UpsertingS3Storages, s3Storage)
}
