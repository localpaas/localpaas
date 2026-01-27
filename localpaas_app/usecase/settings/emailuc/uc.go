package emailuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type EmailUC struct {
	db                       *database.DB
	settingRepo              repository.SettingRepo
	projectSharedSettingRepo repository.ProjectSharedSettingRepo
	permissionManager        permission.Manager
	settingService           settingservice.SettingService
}

func NewEmailUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	permissionManager permission.Manager,
	settingService settingservice.SettingService,
) *EmailUC {
	return &EmailUC{
		db:                       db,
		settingRepo:              settingRepo,
		projectSharedSettingRepo: projectSharedSettingRepo,
		permissionManager:        permissionManager,
		settingService:           settingService,
	}
}
