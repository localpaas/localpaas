package clusterservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type PersistingClusterData struct {
	UpsertingSettings []*entity.Setting
}

func (s *clusterService) PersistClusterData(ctx context.Context, db database.IDB,
	persistingData *PersistingClusterData) error {
	// Persists data
	// Settings
	err := s.settingRepo.UpsertMulti(ctx, db, persistingData.UpsertingSettings,
		entity.SettingUpsertingConflictCols, entity.SettingUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
