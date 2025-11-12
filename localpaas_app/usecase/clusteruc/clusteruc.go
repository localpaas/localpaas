package clusteruc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
	"github.com/localpaas/localpaas/services/docker"
)

type ClusterUC struct {
	db                *database.DB
	permissionManager permission.Manager
	clusterService    clusterservice.ClusterService
	dockerManager     *docker.Manager
}

func NewClusterUC(
	db *database.DB,
	permissionManager permission.Manager,
	clusterService clusterservice.ClusterService,
	dockerManager *docker.Manager,
) *ClusterUC {
	return &ClusterUC{
		db:                db,
		permissionManager: permissionManager,
		clusterService:    clusterService,
		dockerManager:     dockerManager,
	}
}
