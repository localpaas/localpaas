package userserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

func New(
	userRepo repository.UserRepo,
	settingRepo repository.SettingRepo,
	binObjectRepo repository.BinObjectRepo,
	permissionManager permission.Manager,
) userservice.Service {
	return &service{
		userRepo:          userRepo,
		settingRepo:       settingRepo,
		binObjectRepo:     binObjectRepo,
		permissionManager: permissionManager,
	}
}

type service struct {
	userRepo          repository.UserRepo
	settingRepo       repository.SettingRepo
	binObjectRepo     repository.BinObjectRepo
	permissionManager permission.Manager
}
