package cronjobuc

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/cronjobservice"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CronJobUC struct {
	*settings.BaseSettingUC
	appRepo        repository.AppRepo
	taskRepo       repository.TaskRepo
	taskService    taskservice.TaskService
	cronJobService cronjobservice.CronJobService
	taskQueue      queue.TaskQueue
}

func NewCronJobUC(
	baseSettingUC *settings.BaseSettingUC,
	appRepo repository.AppRepo,
	taskRepo repository.TaskRepo,
	taskService taskservice.TaskService,
	cronJobService cronjobservice.CronJobService,
	taskQueue queue.TaskQueue,
) *CronJobUC {
	return &CronJobUC{
		BaseSettingUC:  baseSettingUC,
		appRepo:        appRepo,
		taskRepo:       taskRepo,
		taskService:    taskService,
		cronJobService: cronJobService,
		taskQueue:      taskQueue,
	}
}
