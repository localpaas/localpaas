package appservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type PersistingAppData struct {
	UpsertingApps     []*entity.App
	UpsertingTags     []*entity.AppTag
	UpsertingSettings []*entity.Setting
	UpsertingAccesses []*entity.ACLPermission

	AppsToDeleteTags []string
}

func (s *appService) PersistAppData(ctx context.Context, db database.IDB,
	persistingData *PersistingAppData) error {
	// Deletes all current linked data if configured
	err := s.appTagRepo.DeleteAllByApps(ctx, db, persistingData.AppsToDeleteTags)
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

	// Apps
	err = s.appRepo.UpsertMulti(ctx, db, persistingData.UpsertingApps,
		entity.AppUpsertingConflictCols, entity.AppUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Tags
	err = s.appTagRepo.UpsertMulti(ctx, db, persistingData.UpsertingTags,
		entity.AppTagUpsertingConflictCols, entity.AppTagUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// App accesses
	err = s.permissionManager.UpdateACLPermissions(ctx, db, persistingData.UpsertingAccesses)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
