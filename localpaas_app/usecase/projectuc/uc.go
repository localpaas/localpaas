package projectuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/networkservice"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
)

type ProjectUC struct {
	db                       *database.DB
	projectRepo              repository.ProjectRepo
	settingRepo              repository.SettingRepo
	projectSharedSettingRepo repository.ProjectSharedSettingRepo
	permissionManager        permission.Manager
	userService              userservice.UserService
	projectService           projectservice.ProjectService
	networkService           networkservice.NetworkService
	dockerManager            *docker.Manager
}

func NewProjectUC(
	db *database.DB,
	projectRepo repository.ProjectRepo,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
	projectService projectservice.ProjectService,
	networkService networkservice.NetworkService,
	dockerManager *docker.Manager,
) *ProjectUC {
	return &ProjectUC{
		db:                       db,
		projectRepo:              projectRepo,
		settingRepo:              settingRepo,
		projectSharedSettingRepo: projectSharedSettingRepo,
		permissionManager:        permissionManager,
		userService:              userService,
		projectService:           projectService,
		networkService:           networkService,
		dockerManager:            dockerManager,
	}
}
