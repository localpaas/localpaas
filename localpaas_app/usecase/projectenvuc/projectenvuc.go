package projectenvuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type ProjectEnvUC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	projectRepo       repository.ProjectRepo
	projectEnvRepo    repository.ProjectEnvRepo
	permissionManager permission.Manager
	userService       userservice.UserService
	projectService    projectservice.ProjectService
}

func NewProjectEnvUC(
	db *database.DB,
	userRepo repository.UserRepo,
	projectRepo repository.ProjectRepo,
	projectEnvRepo repository.ProjectEnvRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
	projectService projectservice.ProjectService,
) *ProjectEnvUC {
	return &ProjectEnvUC{
		db:                db,
		userRepo:          userRepo,
		projectRepo:       projectRepo,
		projectEnvRepo:    projectEnvRepo,
		permissionManager: permissionManager,
		userService:       userService,
		projectService:    projectService,
	}
}
