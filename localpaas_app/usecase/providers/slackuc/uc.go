package slackuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type SlackUC struct {
	db                *database.DB
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	settingService    settingservice.SettingService
}

func NewSlackUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	settingService settingservice.SettingService,
) *SlackUC {
	return &SlackUC{
		db:                db,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		settingService:    settingService,
	}
}
