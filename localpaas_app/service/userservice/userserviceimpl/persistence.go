package userserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

func (s *service) PersistUserData(ctx context.Context, db database.IDB,
	persistingData *userservice.PersistingUserData) error {
	// Persists data
	// Users
	err := s.userRepo.UpsertMulti(ctx, db, persistingData.UpsertingUsers,
		entity.UserUpsertingConflictCols, entity.UserUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Settings
	err = s.settingRepo.UpsertMulti(ctx, db, persistingData.UpsertingSettings,
		entity.SettingUpsertingConflictCols, entity.SettingUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Remove accesses
	err = s.permissionManager.RemoveACLPermissions(ctx, db, persistingData.DeletingAccesses)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Project/App/... accesses
	err = s.permissionManager.UpdateACLPermissions(ctx, db, persistingData.UpsertingAccesses)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
