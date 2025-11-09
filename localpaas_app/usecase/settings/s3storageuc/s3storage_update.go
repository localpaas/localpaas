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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
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
		err = uc.prepareUpdatingS3Storage(req.S3StoragePartialReq, s3storageData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.UpdateS3StorageResp{}, nil
}

type updateS3StorageData struct {
	Setting *entity.Setting
}

func (uc *S3StorageUC) loadS3StorageDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *s3storagedto.UpdateS3StorageReq,
	data *updateS3StorageData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, req.ID,
		bunex.SelectFor("UPDATE OF setting"),
		bunex.SelectWhere("setting.type = ?", base.SettingTypeS3Storage),
		bunex.SelectRelation("ObjectAccesses",
			bunex.SelectWhere("acl_permission.subject_type IN (?)", bunex.In([]base.SubjectType{
				base.SubjectTypeProject, base.SubjectTypeApp,
			})),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if req.Name != nil && !strings.EqualFold(setting.Name, *req.Name) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeS3Storage, *req.Name)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("S3Storage").
				WithMsgLog("s3 storage '%s' already exists", conflictSetting.Name)
		}
	}

	return nil
}

func (uc *S3StorageUC) prepareUpdatingS3Storage(
	req *s3storagedto.S3StoragePartialReq,
	data *updateS3StorageData,
	persistingData *persistingS3StorageData,
) error {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	if req.Name != nil {
		setting.Name = *req.Name
	}

	s3Storage, err := setting.ParseS3Storage(false)
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
		s3Storage.SecretKey = *req.SecretKey
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

	err = s3Storage.Encrypt()
	if err != nil {
		return apperrors.Wrap(err)
	}
	setting.MustSetData(s3Storage)

	setting.UpdatedAt = timeNow
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	// Project accesses change
	if req.ProjectAccesses != nil {
		// Remove all current items
		persistingData.DeletingAccesses = append(persistingData.DeletingAccesses, setting.ObjectAccesses...)
		uc.preparePersistingS3StorageProjects(setting, req.ProjectAccesses, timeNow, persistingData)
	}
	return nil
}
