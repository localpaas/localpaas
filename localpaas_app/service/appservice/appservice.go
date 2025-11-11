package appservice

import (
	"context"

	"github.com/localpaas/localpaas/infrastructure/docker"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type AppService interface {
	PersistAppData(ctx context.Context, db database.IDB, data *PersistingAppData) error
}

func NewAppService(
	appRepo repository.AppRepo,
	appTagRepo repository.AppTagRepo,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
	dockerManager *docker.Manager,
) AppService {
	return &appService{
		appRepo:           appRepo,
		appTagRepo:        appTagRepo,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		userService:       userService,
		dockerManager:     dockerManager,
	}
}

type appService struct {
	appRepo           repository.AppRepo
	appTagRepo        repository.AppTagRepo
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	userService       userservice.UserService
	dockerManager     *docker.Manager
}
