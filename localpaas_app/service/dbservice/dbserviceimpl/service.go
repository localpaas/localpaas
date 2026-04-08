package dbserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/dbservice"
)

func New(
	dataMigrationRepo repository.DataMigrationRepo,
	settingRepo repository.SettingRepo,
) dbservice.Service {
	return &service{
		dataMigrationRepo: dataMigrationRepo,
		settingRepo:       settingRepo,
	}
}

type service struct {
	dataMigrationRepo repository.DataMigrationRepo
	settingRepo       repository.SettingRepo
}
