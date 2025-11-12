package clusterservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/services/docker"
)

type ClusterService interface {
	PersistClusterData(ctx context.Context, db database.IDB, data *PersistingClusterData) error
}

func NewClusterService(
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	dockerManager *docker.Manager,
) ClusterService {
	return &clusterService{
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		dockerManager:     dockerManager,
	}
}

type clusterService struct {
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	dockerManager     *docker.Manager
}
