package useruc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type UserUC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	permissionManager permission.Manager
	userService       userservice.UserService
}

func NewUserUC(
	db *database.DB,
	userRepo repository.UserRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
) *UserUC {
	return &UserUC{
		db:                db,
		userRepo:          userRepo,
		permissionManager: permissionManager,
		userService:       userService,
	}
}
