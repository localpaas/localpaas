package volumeuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
	"github.com/localpaas/localpaas/services/docker"
)

type VolumeUC struct {
	db             *database.DB
	settingRepo    repository.SettingRepo
	clusterService clusterservice.ClusterService
	dockerManager  docker.Manager
}

func NewVolumeUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	clusterService clusterservice.ClusterService,
	dockerManager docker.Manager,
) *VolumeUC {
	return &VolumeUC{
		db:             db,
		settingRepo:    settingRepo,
		clusterService: clusterService,
		dockerManager:  dockerManager,
	}
}
