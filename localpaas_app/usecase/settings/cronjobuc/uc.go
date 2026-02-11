package cronjobuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/cronjobservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
)

type CronJobUC struct {
	db                       *database.DB
	settingRepo              repository.SettingRepo
	projectSharedSettingRepo repository.ProjectSharedSettingRepo
	appRepo                  repository.AppRepo
	taskRepo                 repository.TaskRepo
	settingService           settingservice.SettingService
	taskService              taskservice.TaskService
	cronJobService           cronjobservice.CronJobService
	taskQueue                queue.TaskQueue
}

func NewCronJobUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	appRepo repository.AppRepo,
	taskRepo repository.TaskRepo,
	settingService settingservice.SettingService,
	taskService taskservice.TaskService,
	cronJobService cronjobservice.CronJobService,
	taskQueue queue.TaskQueue,
) *CronJobUC {
	return &CronJobUC{
		db:                       db,
		settingRepo:              settingRepo,
		projectSharedSettingRepo: projectSharedSettingRepo,
		appRepo:                  appRepo,
		taskRepo:                 taskRepo,
		settingService:           settingService,
		taskService:              taskService,
		cronJobService:           cronJobService,
		taskQueue:                taskQueue,
	}
}
