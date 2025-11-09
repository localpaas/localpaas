package secretuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type SecretUC struct {
	db                *database.DB
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	settingService    settingservice.SettingService
}

func NewSecretUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	settingService settingservice.SettingService,
) *SecretUC {
	return &SecretUC{
		db:                db,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		settingService:    settingService,
	}
}
