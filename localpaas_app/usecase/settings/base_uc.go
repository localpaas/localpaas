package settings

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/fileservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type BaseSettingUC struct {
	DB                       *database.DB
	SettingRepo              repository.SettingRepo
	ProjectSharedSettingRepo repository.ProjectSharedSettingRepo
	SettingService           settingservice.Service
	FileService              fileservice.Service
}

func NewBaseSettingUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	settingService settingservice.Service,
	fileService fileservice.Service,
) *BaseSettingUC {
	return &BaseSettingUC{
		DB:                       db,
		SettingRepo:              settingRepo,
		ProjectSharedSettingRepo: projectSharedSettingRepo,
		SettingService:           settingService,
		FileService:              fileService,
	}
}
