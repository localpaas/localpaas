package useruc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/emailservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type UC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	binObjectRepo     repository.BinObjectRepo
	userTokenRepo     cacherepository.UserTokenRepo
	permissionManager permission.Manager
	userService       userservice.Service
	emailService      emailservice.Service
}

func New(
	db *database.DB,
	userRepo repository.UserRepo,
	binObjectRepo repository.BinObjectRepo,
	userTokenRepo cacherepository.UserTokenRepo,
	permissionManager permission.Manager,
	userService userservice.Service,
	emailService emailservice.Service,
) *UC {
	return &UC{
		db:                db,
		userRepo:          userRepo,
		binObjectRepo:     binObjectRepo,
		userTokenRepo:     userTokenRepo,
		permissionManager: permissionManager,
		userService:       userService,
		emailService:      emailService,
	}
}
