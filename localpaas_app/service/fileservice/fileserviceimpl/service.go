package fileserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/fileservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

func New(
	settingRepo repository.SettingRepo,
	settingService settingservice.Service,
) fileservice.Service {
	return &service{
		settingRepo:    settingRepo,
		settingService: settingService,
	}
}

type service struct {
	settingRepo    repository.SettingRepo
	settingService settingservice.Service
}
