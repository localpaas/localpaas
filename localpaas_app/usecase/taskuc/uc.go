package taskuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
)

type TaskUC struct {
	db              *database.DB
	taskRepo        repository.TaskRepo
	taskInfoRepo    cacherepository.TaskInfoRepo
	taskControlRepo cacherepository.TaskControlRepo
	taskQueue       taskqueue.TaskQueue
}

func NewTaskUC(
	db *database.DB,
	taskRepo repository.TaskRepo,
	taskInfoRepo cacherepository.TaskInfoRepo,
	taskControlRepo cacherepository.TaskControlRepo,
	taskQueue taskqueue.TaskQueue,
) *TaskUC {
	return &TaskUC{
		db:              db,
		taskRepo:        taskRepo,
		taskInfoRepo:    taskInfoRepo,
		taskControlRepo: taskControlRepo,
		taskQueue:       taskQueue,
	}
}
