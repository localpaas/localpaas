package domainserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/domainservice"
)

func New(
	settingRepo repository.SettingRepo,
	resLinkRepo repository.ResLinkRepo,
) domainservice.Service {
	return &service{
		settingRepo: settingRepo,
		resLinkRepo: resLinkRepo,
	}
}

type service struct {
	settingRepo repository.SettingRepo
	resLinkRepo repository.ResLinkRepo
}
