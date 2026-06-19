package agentserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/agentservice"
	"github.com/localpaas/localpaas/services/docker"
)

func New(
	logger logging.Logger,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	dockerManager docker.Manager,
) agentservice.Service {
	return &service{
		logger:            logger,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		dockerManager:     dockerManager,
	}
}

type service struct {
	logger            logging.Logger
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	dockerManager     docker.Manager
}
