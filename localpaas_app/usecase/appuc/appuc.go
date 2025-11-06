package appuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type AppUC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	projectRepo       repository.ProjectRepo
	appRepo           repository.AppRepo
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	userService       userservice.UserService
	appService        appservice.AppService
}

func NewAppUC(
	db *database.DB,
	userRepo repository.UserRepo,
	projectRepo repository.ProjectRepo,
	appRepo repository.AppRepo,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
	appService appservice.AppService,
) *AppUC {
	return &AppUC{
		db:                db,
		userRepo:          userRepo,
		projectRepo:       projectRepo,
		appRepo:           appRepo,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		userService:       userService,
		appService:        appService,
	}
}
