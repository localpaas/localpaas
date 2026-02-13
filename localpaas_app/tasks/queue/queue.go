package queue

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/gocronqueue"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/cronjobservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
)

type TaskQueue interface {
	Start() error
	Shutdown() error
	RegisterExecutor(typ base.TaskType, execFunc TaskExecFunc)
	RegisterHealthcheckExecutor(execFunc HealthcheckExecFunc)

	ScheduleTask(ctx context.Context, tasks ...*entity.Task) error
	UnscheduleTask(ctx context.Context, tasks ...*entity.Task) error
	ScheduleTasksForCronJob(ctx context.Context, db database.Tx, cronJob *entity.Setting,
		unscheduleCurrentTasks bool) error
}

type taskQueue struct {
	db                        *database.DB
	config                    *config.Config
	logger                    logging.Logger
	server                    *gocronqueue.Server
	client                    *gocronqueue.Client
	redisClient               rediscache.Client
	settingRepo               repository.SettingRepo
	taskRepo                  repository.TaskRepo
	taskInfoRepo              cacherepository.TaskInfoRepo
	healthcheckSettingsRepo   cacherepository.HealthcheckSettingsRepo
	healthcheckNotifEventRepo cacherepository.HealthcheckNotifEventRepo
	cronJobService            cronjobservice.CronJobService
	taskService               taskservice.TaskService
	settingService            settingservice.SettingService

	taskExecutorMap     map[base.TaskType]gocronqueue.TaskExecFunc
	healthcheckExecutor HealthcheckExecFunc
}

func NewTaskQueue(
	db *database.DB,
	config *config.Config,
	logger logging.Logger,
	redisClient rediscache.Client,
	settingRepo repository.SettingRepo,
	taskRepo repository.TaskRepo,
	cacheTaskInfoRepo cacherepository.TaskInfoRepo,
	healthcheckSettingsRepo cacherepository.HealthcheckSettingsRepo,
	healthcheckNotifEventRepo cacherepository.HealthcheckNotifEventRepo,
	cronJobService cronjobservice.CronJobService,
	taskService taskservice.TaskService,
	settingService settingservice.SettingService,
) TaskQueue {
	return &taskQueue{
		db:                        db,
		config:                    config,
		logger:                    logger,
		redisClient:               redisClient,
		settingRepo:               settingRepo,
		taskRepo:                  taskRepo,
		taskInfoRepo:              cacheTaskInfoRepo,
		healthcheckSettingsRepo:   healthcheckSettingsRepo,
		healthcheckNotifEventRepo: healthcheckNotifEventRepo,
		cronJobService:            cronJobService,
		taskService:               taskService,
		settingService:            settingService,
	}
}

func (q *taskQueue) Start() (err error) {
	// Initialize task queue worker if configured
	if q.isWorkerMode() {
		q.logger.Infof("starting task queue server...")
		q.server, err = gocronqueue.NewServer(&gocronqueue.Config{
			TaskMap:                 q.taskExecutorMap,
			RedisClient:             q.redisClient,
			Logger:                  q.logger,
			Concurrency:             q.config.Tasks.Queue.Concurrency,
			TaskCheckInterval:       q.config.Tasks.Queue.TaskCheckInterval,
			TaskCheckFunc:           q.doScheduleTasks,
			TaskCreateInterval:      q.config.Tasks.Queue.TaskCreateInterval,
			TaskCreateFunc:          q.doCreateTasks,
			HealthcheckBaseInterval: q.config.Tasks.Healthcheck.BaseInterval,
			HealthcheckFunc:         q.doHealthcheck,
		})
		if err != nil {
			return apperrors.Wrap(err)
		}

		go func() {
			if err = q.server.Start(); err != nil {
				q.logger.Errorf("failed to start task queue server: %v", err)
			}
		}()
	}

	// Initialize task queue client if configured
	if q.config.RunMode == config.RunModeApp || q.config.RunMode == config.RunModeEmbeddedWorker {
		q.logger.Infof("starting task queue client...")
		q.client, err = gocronqueue.NewClient(q.redisClient, q.logger)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func (q *taskQueue) Shutdown() error {
	q.logger.Info("stopping task queue ...")
	if q.server != nil {
		if err := q.server.Shutdown(); err != nil {
			q.logger.Errorf("failed to start task queue server: %v", err)
			return apperrors.Wrap(err)
		}
	}
	if q.client != nil {
		if err := q.client.Close(); err != nil {
			q.logger.Errorf("failed to stop task queue client: %v", err)
			return apperrors.Wrap(err)
		}
	}
	return nil
}

func (q *taskQueue) isWorkerMode() bool {
	return q.config.RunMode == config.RunModeWorker || q.config.RunMode == config.RunModeEmbeddedWorker
}
