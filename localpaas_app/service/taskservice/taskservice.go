package taskservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
)

type TaskService interface {
	// Logs
	GetTaskLogs(ctx context.Context, db database.IDB, req *GetTaskLogsReq) (*GetTaskLogsResp, error)
}

func NewTaskService(
	redisClient rediscache.Client,
	taskRepo repository.TaskRepo,
	taskLogRepo repository.TaskLogRepo,
	settingRepo repository.SettingRepo,
	taskInfoRepo cacherepository.TaskInfoRepo,
) TaskService {
	return &taskService{
		redisClient:  redisClient,
		taskRepo:     taskRepo,
		taskLogRepo:  taskLogRepo,
		settingRepo:  settingRepo,
		taskInfoRepo: taskInfoRepo,
	}
}

type taskService struct {
	redisClient  rediscache.Client
	taskRepo     repository.TaskRepo
	taskLogRepo  repository.TaskLogRepo
	settingRepo  repository.SettingRepo
	taskInfoRepo cacherepository.TaskInfoRepo
}
