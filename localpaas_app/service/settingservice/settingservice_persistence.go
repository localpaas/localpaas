package settingservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type PersistingSettingData struct {
	UpsertingSettings   []*entity.Setting
	UpsertingS3Storages []*entity.S3Storage
	UpsertingAccesses   []*entity.ACLPermission
	DeletingAccesses    []*entity.ACLPermission
}

func (s *settingService) PersistSettingData(ctx context.Context, db database.IDB,
	persistingData *PersistingSettingData) error {
	// Deletes data
	err := s.permissionManager.UpdateACLPermissions(ctx, db, persistingData.DeletingAccesses)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Persists data
	// Settings
	err = s.settingRepo.UpsertMulti(ctx, db, persistingData.UpsertingSettings,
		entity.SettingUpsertingConflictCols, entity.SettingUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// S3 storages
	err = s.s3StorageRepo.UpsertMulti(ctx, db, persistingData.UpsertingS3Storages,
		entity.S3StorageUpsertingConflictCols, entity.S3StorageUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Accesses
	err = s.permissionManager.UpdateACLPermissions(ctx, db, persistingData.UpsertingAccesses)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
