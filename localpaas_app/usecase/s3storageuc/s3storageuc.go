package s3storageuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type S3StorageUC struct {
	db                *database.DB
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	settingService    settingservice.SettingService
}

func NewS3StorageUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	settingService settingservice.SettingService,
) *S3StorageUC {
	return &S3StorageUC{
		db:                db,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		settingService:    settingService,
	}
}
