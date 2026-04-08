package cronjobserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/cronjobservice"
	"github.com/localpaas/localpaas/localpaas_app/service/envvarservice"
)

func New(
	settingRepo repository.SettingRepo,
	envVarService envvarservice.Service,
) cronjobservice.Service {
	return &service{
		settingRepo:   settingRepo,
		envVarService: envVarService,
	}
}

type service struct {
	settingRepo   repository.SettingRepo
	envVarService envvarservice.Service
}
