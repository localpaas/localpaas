package imageuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
	"github.com/localpaas/localpaas/services/docker"
)

type ImageUC struct {
	db                *database.DB
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	clusterService    clusterservice.ClusterService
	dockerManager     *docker.Manager
}

func NewImageUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	clusterService clusterservice.ClusterService,
	dockerManager *docker.Manager,
) *ImageUC {
	return &ImageUC{
		db:                db,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		clusterService:    clusterService,
		dockerManager:     dockerManager,
	}
}
