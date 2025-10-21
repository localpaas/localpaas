package projectservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type ProjectService interface {
	PersistProjectData(ctx context.Context, db database.IDB, data *PersistingProjectData) error
}

func NewProjectService(
	projectRepo repository.ProjectRepo,
	projectEnvRepo repository.ProjectEnvRepo,
	projectTagRepo repository.ProjectTagRepo,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
) ProjectService {
	return &projectService{
		projectRepo:       projectRepo,
		projectEnvRepo:    projectEnvRepo,
		projectTagRepo:    projectTagRepo,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		userService:       userService,
	}
}

type projectService struct {
	projectRepo       repository.ProjectRepo
	projectEnvRepo    repository.ProjectEnvRepo
	projectTagRepo    repository.ProjectTagRepo
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	userService       userservice.UserService
}
