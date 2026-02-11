package healthcheckuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
)

type HealthcheckUC struct {
	db                       *database.DB
	settingRepo              repository.SettingRepo
	projectSharedSettingRepo repository.ProjectSharedSettingRepo
	appRepo                  repository.AppRepo
	taskRepo                 repository.TaskRepo
	settingService           settingservice.SettingService
	taskService              taskservice.TaskService
	taskQueue                queue.TaskQueue
}

func NewHealthcheckUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	appRepo repository.AppRepo,
	taskRepo repository.TaskRepo,
	settingService settingservice.SettingService,
	taskService taskservice.TaskService,
	taskQueue queue.TaskQueue,
) *HealthcheckUC {
	return &HealthcheckUC{
		db:                       db,
		settingRepo:              settingRepo,
		projectSharedSettingRepo: projectSharedSettingRepo,
		appRepo:                  appRepo,
		taskRepo:                 taskRepo,
		settingService:           settingService,
		taskService:              taskService,
		taskQueue:                taskQueue,
	}
}
