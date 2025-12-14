package cronjobuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
)

type CronJobUC struct {
	db             *database.DB
	settingRepo    repository.SettingRepo
	taskRepo       repository.TaskRepo
	settingService settingservice.SettingService
	taskQueue      taskqueue.TaskQueue
}

func NewCronJobUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	taskRepo repository.TaskRepo,
	settingService settingservice.SettingService,
	taskQueue taskqueue.TaskQueue,
) *CronJobUC {
	return &CronJobUC{
		db:             db,
		settingRepo:    settingRepo,
		taskRepo:       taskRepo,
		settingService: settingService,
		taskQueue:      taskQueue,
	}
}
