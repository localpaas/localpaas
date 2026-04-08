package traefikserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/traefikservice"
	"github.com/localpaas/localpaas/services/docker"
)

func New(
	settingRepo repository.SettingRepo,
	dockerManager docker.Manager,
) traefikservice.Service {
	return &service{
		settingRepo:   settingRepo,
		dockerManager: dockerManager,
	}
}

type service struct {
	settingRepo   repository.SettingRepo
	dockerManager docker.Manager
}
