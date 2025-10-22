package clusteruc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
)

type ClusterUC struct {
	db                *database.DB
	nodeRepo          repository.NodeRepo
	permissionManager permission.Manager
	clusterService    clusterservice.ClusterService
}

func NewClusterUC(
	db *database.DB,
	nodeRepo repository.NodeRepo,
	permissionManager permission.Manager,
	clusterService clusterservice.ClusterService,
) *ClusterUC {
	return &ClusterUC{
		db:                db,
		nodeRepo:          nodeRepo,
		permissionManager: permissionManager,
		clusterService:    clusterService,
	}
}
