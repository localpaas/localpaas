package projectuc

import (
	"github.com/localpaas/localpaas/infrastructure/docker"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type ProjectUC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	projectRepo       repository.ProjectRepo
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	userService       userservice.UserService
	projectService    projectservice.ProjectService
	dockerManager     *docker.Manager
}

func NewProjectUC(
	db *database.DB,
	userRepo repository.UserRepo,
	projectRepo repository.ProjectRepo,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
	projectService projectservice.ProjectService,
	dockerManager *docker.Manager,
) *ProjectUC {
	return &ProjectUC{
		db:                db,
		userRepo:          userRepo,
		projectRepo:       projectRepo,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		userService:       userService,
		projectService:    projectService,
		dockerManager:     dockerManager,
	}
}
