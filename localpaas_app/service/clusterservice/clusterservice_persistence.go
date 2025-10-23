package clusterservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type PersistingClusterData struct {
	UpsertingNodes    []*entity.Node
	UpsertingSettings []*entity.Setting

	NodesToDeleteSettings []string
}

func (s *clusterService) PersistClusterData(ctx context.Context, db database.IDB,
	persistingData *PersistingClusterData) error {
	// Deletes all current linked data if configured
	// Main settings
	err := s.settingRepo.DeleteAllByTargetObjects(ctx, db, base.SettingTypeNode,
		persistingData.NodesToDeleteSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Persists data
	// Nodes
	err = s.nodeRepo.UpsertMulti(ctx, db, persistingData.UpsertingNodes,
		entity.NodeUpsertingConflictCols, entity.NodeUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Settings
	err = s.settingRepo.UpsertMulti(ctx, db, persistingData.UpsertingSettings,
		entity.SettingUpsertingConflictCols, entity.SettingUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
