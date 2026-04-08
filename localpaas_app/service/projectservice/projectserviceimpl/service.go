package projectserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/networkservice"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
)

func New(
	projectRepo repository.ProjectRepo,
	appRepo repository.AppRepo,
	projectTagRepo repository.ProjectTagRepo,
	settingRepo repository.SettingRepo,
	userRepo repository.UserRepo,
	permissionManager permission.Manager,
	userService userservice.Service,
	appService appservice.Service,
	networkService networkservice.Service,
	dockerManager docker.Manager,
) projectservice.Service {
	return &service{
		projectRepo:       projectRepo,
		appRepo:           appRepo,
		projectTagRepo:    projectTagRepo,
		settingRepo:       settingRepo,
		userRepo:          userRepo,
		permissionManager: permissionManager,
		userService:       userService,
		appService:        appService,
		networkService:    networkService,
		dockerManager:     dockerManager,
	}
}

type service struct {
	projectRepo       repository.ProjectRepo
	appRepo           repository.AppRepo
	projectTagRepo    repository.ProjectTagRepo
	settingRepo       repository.SettingRepo
	userRepo          repository.UserRepo
	permissionManager permission.Manager
	userService       userservice.Service
	appService        appservice.Service
	networkService    networkservice.Service
	dockerManager     docker.Manager
}
