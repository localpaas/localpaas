package systembackupuc

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/cronjobservice"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type SystemBackupUC struct {
	*settings.BaseSettingUC
	appRepo        repository.AppRepo
	taskRepo       repository.TaskRepo
	taskService    taskservice.Service
	cronJobService cronjobservice.Service
	taskQueue      queue.TaskQueue
}

func NewSystemBackupUC(
	baseSettingUC *settings.BaseSettingUC,
	appRepo repository.AppRepo,
	taskRepo repository.TaskRepo,
	taskService taskservice.Service,
	cronJobService cronjobservice.Service,
	taskQueue queue.TaskQueue,
) *SystemBackupUC {
	return &SystemBackupUC{
		BaseSettingUC:  baseSettingUC,
		appRepo:        appRepo,
		taskRepo:       taskRepo,
		taskService:    taskService,
		cronJobService: cronJobService,
		taskQueue:      taskQueue,
	}
}
