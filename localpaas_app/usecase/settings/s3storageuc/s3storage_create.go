package s3storageuc

import (
	"context"
	"errors"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

func (uc *S3StorageUC) CreateS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.CreateS3StorageReq,
) (*s3storagedto.CreateS3StorageResp, error) {
	s3storageData := &createS3StorageData{}
	err := uc.loadS3StorageData(ctx, uc.db, req, s3storageData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingS3StorageData{}
	uc.preparePersistingS3Storage(req.S3StorageBaseReq, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &s3storagedto.CreateS3StorageResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createS3StorageData struct {
}

func (uc *S3StorageUC) loadS3StorageData(
	ctx context.Context,
	db database.IDB,
	req *s3storagedto.CreateS3StorageReq,
	_ *createS3StorageData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeS3Storage, req.Name)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("S3Storage").
			WithMsgLog("s3 storage '%s' setting already exists", req.Name)
	}

	return nil
}

type persistingS3StorageData struct {
	settingservice.PersistingSettingData
}

func (uc *S3StorageUC) preparePersistingS3Storage(
	req *s3storagedto.S3StorageBaseReq,
	persistingData *persistingS3StorageData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeS3Storage,
		Status:    base.SettingStatusActive,
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	s3Storage := &entity.S3Storage{
		AccessKeyID: req.AccessKeyID,
		SecretKey:   req.SecretKey,
		Region:      req.Region,
		Bucket:      req.Bucket,
		Endpoint:    req.Endpoint,
	}
	setting.MustSetData(s3Storage.MustEncrypt())

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	uc.preparePersistingS3StorageProjects(setting, req.ProjectAccesses, timeNow, persistingData)
}

func (uc *S3StorageUC) preparePersistingS3StorageProjects(
	setting *entity.Setting,
	projectReqs []*s3storagedto.S3StorageProjectAccessReq,
	timeNow time.Time,
	persistingData *persistingS3StorageData,
) {
	for _, projectReq := range projectReqs {
		persistingData.UpsertingAccesses = append(persistingData.UpsertingAccesses,
			&entity.ACLPermission{
				SubjectType:  base.SubjectTypeProject,
				SubjectID:    projectReq.ID,
				ResourceType: base.ResourceTypeS3Storage,
				ResourceID:   setting.ID,
				Actions:      entity.AccessActions{Read: projectReq.Allowed},
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			})
		uc.preparePersistingS3StorageApps(setting, projectReq.AppAccesses, timeNow, persistingData)
	}
}

func (uc *S3StorageUC) preparePersistingS3StorageApps(
	setting *entity.Setting,
	appReqs []*s3storagedto.S3StorageAppAccessReq,
	timeNow time.Time,
	persistingData *persistingS3StorageData,
) {
	for _, appReq := range appReqs {
		persistingData.UpsertingAccesses = append(persistingData.UpsertingAccesses,
			&entity.ACLPermission{
				SubjectType:  base.SubjectTypeApp,
				SubjectID:    appReq.ID,
				ResourceType: base.ResourceTypeS3Storage,
				ResourceID:   setting.ID,
				Actions:      entity.AccessActions{Read: appReq.Allowed},
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			})
	}
}

func (uc *S3StorageUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingS3StorageData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
