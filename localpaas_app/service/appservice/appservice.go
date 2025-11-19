package appservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
)

type AppService interface {
	PersistAppData(ctx context.Context, db database.IDB, data *PersistingAppData) error

	UpdateAppDeployment(ctx context.Context, app *entity.App, req *AppDeploymentReq) (
		*AppDeploymentResp, error)
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
