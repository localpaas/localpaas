package cronjobuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type CronJobUC struct {
	db             *database.DB
	settingRepo    repository.SettingRepo
	settingService settingservice.SettingService
}

func NewCronJobUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	settingService settingservice.SettingService,
) *CronJobUC {
	return &CronJobUC{
		db:             db,
		settingRepo:    settingRepo,
		settingService: settingService,
	}
}
