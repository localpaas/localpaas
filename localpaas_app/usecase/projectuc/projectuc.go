package projectuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type ProjectUC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	projectRepo       repository.ProjectRepo
	permissionManager permission.Manager
	userService       userservice.UserService
}

func NewProjectUC(
	db *database.DB,
	userRepo repository.UserRepo,
	projectRepo repository.ProjectRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
) *ProjectUC {
	return &ProjectUC{
		db:                db,
		userRepo:          userRepo,
		projectRepo:       projectRepo,
		permissionManager: permissionManager,
		userService:       userService,
	}
}
