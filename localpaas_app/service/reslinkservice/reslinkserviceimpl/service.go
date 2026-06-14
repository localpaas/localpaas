package reslinkserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/reslinkservice"
)

func New(
	settingRepo repository.SettingRepo,
	resLinkRepo repository.ResLinkRepo,
) reslinkservice.Service {
	return &service{
		settingRepo: settingRepo,
		resLinkRepo: resLinkRepo,
	}
}

type service struct {
	settingRepo repository.SettingRepo
	resLinkRepo repository.ResLinkRepo
}
