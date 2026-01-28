package imserviceuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type IMServiceUC struct {
	db                       *database.DB
	settingRepo              repository.SettingRepo
	projectSharedSettingRepo repository.ProjectSharedSettingRepo
	permissionManager        permission.Manager
	settingService           settingservice.SettingService
}

func NewIMServiceUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	permissionManager permission.Manager,
	settingService settingservice.SettingService,
) *IMServiceUC {
	return &IMServiceUC{
		db:                       db,
		settingRepo:              settingRepo,
		projectSharedSettingRepo: projectSharedSettingRepo,
		permissionManager:        permissionManager,
		settingService:           settingService,
	}
}
