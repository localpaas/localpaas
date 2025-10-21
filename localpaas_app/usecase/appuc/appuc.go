package appuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type AppUC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	appRepo           repository.AppRepo
	permissionManager permission.Manager
	userService       userservice.UserService
}

func NewAppUC(
	db *database.DB,
	userRepo repository.UserRepo,
	appRepo repository.AppRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
) *AppUC {
	return &AppUC{
		db:                db,
		userRepo:          userRepo,
		appRepo:           appRepo,
		permissionManager: permissionManager,
		userService:       userService,
	}
}
