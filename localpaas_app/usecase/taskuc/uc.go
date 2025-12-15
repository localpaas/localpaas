package taskuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
)

type TaskUC struct {
	db                *database.DB
	taskRepo          repository.TaskRepo
	cacheTaskInfoRepo cacherepository.TaskInfoRepo
	taskQueue         taskqueue.TaskQueue
}

func NewTaskUC(
	db *database.DB,
	taskRepo repository.TaskRepo,
	cacheTaskInfoRepo cacherepository.TaskInfoRepo,
	taskQueue taskqueue.TaskQueue,
) *TaskUC {
	return &TaskUC{
		db:                db,
		taskRepo:          taskRepo,
		cacheTaskInfoRepo: cacheTaskInfoRepo,
		taskQueue:         taskQueue,
	}
}
