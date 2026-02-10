package taskservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
)

type TaskService interface {
	GetTask(ctx context.Context, db database.IDB, req *GetTaskReq, extraOpts ...bunex.SelectQueryOption) (
		*GetTaskResp, error)
	ListTask(ctx context.Context, db database.IDB, req *ListTaskReq, extraOpts ...bunex.SelectQueryOption) (
		*ListTaskResp, error)

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
