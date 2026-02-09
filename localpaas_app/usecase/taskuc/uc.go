package taskuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
)

type TaskUC struct {
	db              *database.DB
	redisClient     rediscache.Client
	taskRepo        repository.TaskRepo
	taskLogRepo     repository.TaskLogRepo
	taskInfoRepo    cacherepository.TaskInfoRepo
	taskControlRepo cacherepository.TaskControlRepo
	taskService     taskservice.TaskService
}

func NewTaskUC(
	db *database.DB,
	redisClient rediscache.Client,
	taskRepo repository.TaskRepo,
	taskLogRepo repository.TaskLogRepo,
	taskInfoRepo cacherepository.TaskInfoRepo,
	taskControlRepo cacherepository.TaskControlRepo,
	taskService taskservice.TaskService,
) *TaskUC {
	return &TaskUC{
		db:              db,
		redisClient:     redisClient,
		taskRepo:        taskRepo,
		taskLogRepo:     taskLogRepo,
		taskInfoRepo:    taskInfoRepo,
		taskControlRepo: taskControlRepo,
		taskService:     taskService,
	}
}
