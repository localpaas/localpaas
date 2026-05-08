package appdeploymentserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/appdeploymentservice"
	"github.com/localpaas/localpaas/localpaas_app/service/envvarservice"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
)

type service struct {
	logger              logging.Logger
	db                  *database.DB
	redisClient         rediscache.Client
	settingRepo         repository.SettingRepo
	deploymentRepo      repository.DeploymentRepo
	taskLogRepo         repository.TaskLogRepo
	taskRepo            repository.TaskRepo
	taskInfoRepo        cacherepository.TaskInfoRepo
	deploymentInfoRepo  cacherepository.DeploymentInfoRepo
	dockerManager       docker.Manager
	envVarService       envvarservice.Service
	settingService      settingservice.Service
	userService         userservice.Service
	notificationService notificationservice.Service
}

func New(
	logger logging.Logger,
	db *database.DB,
	redisClient rediscache.Client,
	settingRepo repository.SettingRepo,
	deploymentRepo repository.DeploymentRepo,
	taskLogRepo repository.TaskLogRepo,
	taskRepo repository.TaskRepo,
	taskInfoRepo cacherepository.TaskInfoRepo,
	deploymentInfoRepo cacherepository.DeploymentInfoRepo,
	dockerManager docker.Manager,
	envVarService envvarservice.Service,
	settingService settingservice.Service,
	userService userservice.Service,
	notificationService notificationservice.Service,
) appdeploymentservice.Service {
	s := &service{
		logger:              logger,
		db:                  db,
		redisClient:         redisClient,
		settingRepo:         settingRepo,
		deploymentRepo:      deploymentRepo,
		taskLogRepo:         taskLogRepo,
		taskRepo:            taskRepo,
		taskInfoRepo:        taskInfoRepo,
		deploymentInfoRepo:  deploymentInfoRepo,
		dockerManager:       dockerManager,
		envVarService:       envVarService,
		settingService:      settingService,
		userService:         userService,
		notificationService: notificationService,
	}
	return s
}
