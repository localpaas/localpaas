package s3storageuc

import (
	"context"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/s3storageuc/s3storagedto"
	"github.com/localpaas/localpaas/pkg/reflectutil"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

func (uc *S3StorageUC) UpdateS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.UpdateS3StorageReq,
) (*s3storagedto.UpdateS3StorageResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		s3storageData := &updateS3StorageData{}
		err := uc.loadS3StorageDataForUpdate(ctx, db, req, s3storageData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingS3StorageData{}
		uc.prepareUpdatingS3Storage(req.S3StoragePartialReq, s3storageData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.UpdateS3StorageResp{}, nil
}

type updateS3StorageData struct {
	S3Storage *entity.S3Storage
}

func (uc *S3StorageUC) loadS3StorageDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *s3storagedto.UpdateS3StorageReq,
	data *updateS3StorageData,
) error {
	s3Storage, err := uc.s3StorageRepo.GetByID(ctx, db, req.ID,
		bunex.SelectFor("UPDATE OF s3_storage"),
		bunex.SelectRelation("ObjectAccesses",
			bunex.SelectWhere("acl_permission.subject_type IN (?)", bunex.In([]base.SubjectType{
				base.SubjectTypeProject, base.SubjectTypeApp,
			})),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.S3Storage = s3Storage

	// If name changes, validate the new one
	if req.Name != nil && !strings.EqualFold(s3Storage.Name, *req.Name) {
		conflictS3Storage, _ := uc.s3StorageRepo.GetByName(ctx, db, *req.Name)
		if conflictS3Storage != nil {
			return apperrors.NewAlreadyExist("S3Storage").
				WithMsgLog("s3 storage '%s' already exists", *req.Name)
		}
	}

	return nil
}

func (uc *S3StorageUC) prepareUpdatingS3Storage(
	req *s3storagedto.S3StoragePartialReq,
	data *updateS3StorageData,
	persistingData *persistingS3StorageData,
) {
	timeNow := timeutil.NowUTC()
	s3Storage := data.S3Storage
	if req.Name != nil {
		s3Storage.Name = *req.Name
	}
	if req.AccessKeyID != nil {
		s3Storage.AccessKeyID = *req.AccessKeyID
	}
	// TODO: encrypt the data (secret access key)
	if req.SecretAccessKey != nil {
		s3Storage.SecretAccessKey = reflectutil.UnsafeStrToBytes(*req.SecretAccessKey)
	}
	if req.Region != nil {
		s3Storage.Region = *req.Region
	}
	if req.Bucket != nil {
		s3Storage.Bucket = *req.Bucket
	}

	persistingData.UpsertingS3Storages = append(persistingData.UpsertingS3Storages, s3Storage)

	// Project accesses change
	if req.ProjectAccesses != nil {
		// Remove all current items
		persistingData.DeletingAccesses = append(persistingData.DeletingAccesses, s3Storage.ObjectAccesses...)
		uc.preparePersistingS3StorageProjects(s3Storage, req.ProjectAccesses, timeNow, persistingData)
	}
}
