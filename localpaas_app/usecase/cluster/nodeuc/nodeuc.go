package nodeuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
	"github.com/localpaas/localpaas/localpaas_app/service/lpappservice"
	"github.com/localpaas/localpaas/services/docker"
)

type NodeUC struct {
	db                *database.DB
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	clusterService    clusterservice.ClusterService
	lpAppService      lpappservice.LpAppService
	dockerManager     *docker.Manager
}

func NewNodeUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	clusterService clusterservice.ClusterService,
	lpAppService lpappservice.LpAppService,
	dockerManager *docker.Manager,
) *NodeUC {
	return &NodeUC{
		db:                db,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		clusterService:    clusterService,
		lpAppService:      lpAppService,
		dockerManager:     dockerManager,
	}
}
