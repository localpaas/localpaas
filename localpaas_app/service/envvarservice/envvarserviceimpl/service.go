package envvarserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/envvarservice"
)

func New(
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
) envvarservice.Service {
	return &service{
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
	}
}

type service struct {
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
}
