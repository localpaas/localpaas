package notificationserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
)

func New(
	settingRepo repository.SettingRepo,
) notificationservice.Service {
	return &service{
		settingRepo: settingRepo,
	}
}

type service struct {
	settingRepo repository.SettingRepo
}
