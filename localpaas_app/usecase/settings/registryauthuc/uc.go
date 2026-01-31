package registryauthuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/services/docker"
)

type RegistryAuthUC struct {
	db                       *database.DB
	settingRepo              repository.SettingRepo
	projectSharedSettingRepo repository.ProjectSharedSettingRepo
	settingService           settingservice.SettingService
	dockerManager            *docker.Manager
}

func NewRegistryAuthUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	settingService settingservice.SettingService,
	dockerManager *docker.Manager,
) *RegistryAuthUC {
	return &RegistryAuthUC{
		db:                       db,
		settingRepo:              settingRepo,
		projectSharedSettingRepo: projectSharedSettingRepo,
		settingService:           settingService,
		dockerManager:            dockerManager,
	}
}
