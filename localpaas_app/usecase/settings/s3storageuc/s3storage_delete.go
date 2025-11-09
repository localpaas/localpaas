package s3storageuc

import (
	"context"

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
	Setting *entity.Setting
}

func (uc *S3StorageUC) loadS3StorageDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *s3storagedto.DeleteS3StorageReq,
	data *deleteS3StorageData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, req.ID,
		bunex.SelectFor("UPDATE OF setting"),
		bunex.SelectWhere("setting.type = ?", base.SettingTypeS3Storage),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *S3StorageUC) prepareDeletingS3Storage(
	data *deleteS3StorageData,
	persistingData *persistingS3StorageData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
