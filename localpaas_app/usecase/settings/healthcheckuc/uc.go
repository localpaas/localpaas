package healthcheckuc

import (
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type HealthcheckUC struct {
	*settings.BaseSettingUC
	appRepo     repository.AppRepo
	taskRepo    repository.TaskRepo
	taskService taskservice.TaskService
	taskQueue   queue.TaskQueue
}

func NewHealthcheckUC(
	baseSettingUC *settings.BaseSettingUC,
	appRepo repository.AppRepo,
	taskRepo repository.TaskRepo,
	taskService taskservice.TaskService,
	taskQueue queue.TaskQueue,
) *HealthcheckUC {
	return &HealthcheckUC{
		BaseSettingUC: baseSettingUC,
		appRepo:       appRepo,
		taskRepo:      taskRepo,
		taskService:   taskService,
		taskQueue:     taskQueue,
	}
}
