package settingserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/sslservice"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
)

func New(
	db *database.DB,
	settingRepo repository.SettingRepo,
	appRepo repository.AppRepo,
	healthcheckSettingsRepo cacherepository.HealthcheckSettingsRepo,
	userService userservice.Service,
	taskService taskservice.Service,
	sslService sslservice.Service,
	permissionManager permission.Manager,
	dockerManager docker.Manager,
) settingservice.Service {
	return &service{
		db:                      db,
		settingRepo:             settingRepo,
		appRepo:                 appRepo,
		healthcheckSettingsRepo: healthcheckSettingsRepo,
		userService:             userService,
		taskService:             taskService,
		sslService:              sslService,
		permissionManager:       permissionManager,
		dockerManager:           dockerManager,
	}
}

type service struct {
	db                      *database.DB
	settingRepo             repository.SettingRepo
	appRepo                 repository.AppRepo
	healthcheckSettingsRepo cacherepository.HealthcheckSettingsRepo
	userService             userservice.Service
	taskService             taskservice.Service
	sslService              sslservice.Service
	permissionManager       permission.Manager
	dockerManager           docker.Manager
}
