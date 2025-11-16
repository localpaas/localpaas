package ssluc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type SslUC struct {
	db                *database.DB
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	settingService    settingservice.SettingService
}

func NewSslUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	settingService settingservice.SettingService,
) *SslUC {
	return &SslUC{
		db:                db,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		settingService:    settingService,
	}
}
