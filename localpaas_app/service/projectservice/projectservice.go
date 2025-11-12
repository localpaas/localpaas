package projectservice

import (
	"context"

	"github.com/docker/docker/api/types/network"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
)

type ProjectService interface {
	PersistProjectData(ctx context.Context, db database.IDB, data *PersistingProjectData) error

	CreateProjectNetworks(ctx context.Context, project *entity.Project) (*network.CreateResponse, error)
	ListProjectNetworks(ctx context.Context, project *entity.Project) ([]network.Summary, error)
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
