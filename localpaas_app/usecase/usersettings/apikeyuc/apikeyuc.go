package apikeyuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type APIKeyUC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	settingService    settingservice.SettingService
}

func NewAPIKeyUC(
	db *database.DB,
	userRepo repository.UserRepo,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	settingService settingservice.SettingService,
) *APIKeyUC {
	return &APIKeyUC{
		db:                db,
		userRepo:          userRepo,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		settingService:    settingService,
	}
}
