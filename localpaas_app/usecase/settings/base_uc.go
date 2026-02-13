package settings

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type BaseSettingUC struct {
	DB                       *database.DB
	SettingRepo              repository.SettingRepo
	ProjectSharedSettingRepo repository.ProjectSharedSettingRepo
	SettingService           settingservice.SettingService
}

func NewBaseSettingUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	settingService settingservice.SettingService,
) *BaseSettingUC {
	return &BaseSettingUC{
		DB:                       db,
		SettingRepo:              settingRepo,
		ProjectSharedSettingRepo: projectSharedSettingRepo,
		SettingService:           settingService,
	}
}
