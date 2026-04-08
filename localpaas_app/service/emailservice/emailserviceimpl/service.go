package emailserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/emailservice"
)

func New(
	settingRepo repository.SettingRepo,
) emailservice.Service {
	return &service{
		settingRepo: settingRepo,
	}
}

type service struct {
	settingRepo repository.SettingRepo
}
