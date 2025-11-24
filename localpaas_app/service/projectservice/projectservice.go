package projectservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
)

type ProjectService interface {
	PersistProjectData(ctx context.Context, db database.IDB, data *PersistingProjectData) error
}

func NewProjectService(
	projectRepo repository.ProjectRepo,
	projectTagRepo repository.ProjectTagRepo,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
	dockerManager *docker.Manager,
) ProjectService {
	return &projectService{
		projectRepo:       projectRepo,
		projectTagRepo:    projectTagRepo,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		userService:       userService,
		dockerManager:     dockerManager,
	}
}

type projectService struct {
	projectRepo       repository.ProjectRepo
	projectTagRepo    repository.ProjectTagRepo
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	userService       userservice.UserService
	dockerManager     *docker.Manager
}
