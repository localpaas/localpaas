package userservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type PersistingUserData struct {
	UpsertingUsers    []*entity.User
	UpsertingSettings []*entity.Setting
	UpsertingAccesses []*entity.ACLPermission
}

func (s *userService) PersistUserData(ctx context.Context, db database.IDB,
	persistingData *PersistingUserData) error {
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

	// App accesses
	err = s.permissionManager.UpdateACLPermissions(ctx, db, persistingData.UpsertingAccesses)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
