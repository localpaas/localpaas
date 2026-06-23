package clusterserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
)

func (s *service) PersistClusterData(ctx context.Context, db database.IDB,
	persistingData *clusterservice.PersistingClusterData) error {
	// Persists data
	// Settings
	err := s.settingRepo.UpsertMulti(ctx, db, persistingData.UpsertingSettings,
		entity.SettingUpsertingConflictCols, entity.SettingUpsertingUpdateCols)
	if err != nil {
		return apperrors.New(err)
	}

	return nil
}
