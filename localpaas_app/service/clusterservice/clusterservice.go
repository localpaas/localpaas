package clusterservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type ClusterService interface {
	PersistClusterData(ctx context.Context, db database.IDB, data *PersistingClusterData) error
}

func NewClusterService(
	nodeRepo repository.NodeRepo,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
) ClusterService {
	return &clusterService{
		nodeRepo:          nodeRepo,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
	}
}

type clusterService struct {
	nodeRepo          repository.NodeRepo
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
}
