package taskservice

import (
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/entity"
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

	// Locking
	CreateDBLock(ctx context.Context, db database.Tx, id, selectFor string) (*entity.Lock, error)
	CreateRedisLock(ctx context.Context, key string, exp time.Duration) (success bool, releaser func(), err error)
}

func NewTaskService(
	db *database.DB,
	redisClient rediscache.Client,
	taskRepo repository.TaskRepo,
	taskLogRepo repository.TaskLogRepo,
	settingRepo repository.SettingRepo,
	lockRepo repository.LockRepo,
	taskInfoRepo cacherepository.TaskInfoRepo,
) TaskService {
	return &taskService{
		db:           db,
		redisClient:  redisClient,
		taskRepo:     taskRepo,
		taskLogRepo:  taskLogRepo,
		settingRepo:  settingRepo,
		lockRepo:     lockRepo,
		taskInfoRepo: taskInfoRepo,
	}
}

type taskService struct {
	db           *database.DB
	redisClient  rediscache.Client
	taskRepo     repository.TaskRepo
	taskLogRepo  repository.TaskLogRepo
	settingRepo  repository.SettingRepo
	lockRepo     repository.LockRepo
	taskInfoRepo cacherepository.TaskInfoRepo
}
