package volumeuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/services/docker"
)

type VolumeUC struct {
	db             *database.DB
	settingRepo    repository.SettingRepo
	clusterService clusterservice.Service
	projectService projectservice.Service
	dockerManager  docker.Manager
}

func NewVolumeUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	clusterService clusterservice.Service,
	projectService projectservice.Service,
	dockerManager docker.Manager,
) *VolumeUC {
	return &VolumeUC{
		db:             db,
		settingRepo:    settingRepo,
		clusterService: clusterService,
		projectService: projectService,
		dockerManager:  dockerManager,
	}
}
