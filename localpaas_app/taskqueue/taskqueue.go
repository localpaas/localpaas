package taskqueue

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/gocronqueue"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type TaskQueue interface {
	Start(cfg *config.Config) error
	Shutdown() error

	ScheduleTasks(ctx context.Context, tasks []*entity.Task) error
	UnscheduleTasks(ctx context.Context, tasks []*entity.Task) error
}

type taskQueue struct {
	db               *database.DB
	logger           logging.Logger
	server           *gocronqueue.Server
	client           *gocronqueue.Client
	settingRepo      repository.SettingRepo
	taskRepo         repository.TaskRepo
	updatingTaskRepo repository.UpdatingTaskRepo
}

func NewTaskQueue(
	db *database.DB,
	logger logging.Logger,
	server *gocronqueue.Server,
	client *gocronqueue.Client,
	settingRepo repository.SettingRepo,
	taskRepo repository.TaskRepo,
	updatingTaskRepo repository.UpdatingTaskRepo,
) TaskQueue {
	return &taskQueue{
		db:               db,
		logger:           logger,
		server:           server,
		client:           client,
		settingRepo:      settingRepo,
		taskRepo:         taskRepo,
		updatingTaskRepo: updatingTaskRepo,
	}
}
