package appdeploymentuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/services/docker"
)

type AppDeploymentUC struct {
	db                 *database.DB
	redisClient        rediscache.Client
	projectRepo        repository.ProjectRepo
	appRepo            repository.AppRepo
	deploymentRepo     repository.DeploymentRepo
	deploymentLogRepo  repository.DeploymentLogRepo
	deploymentInfoRepo cacherepository.DeploymentInfoRepo
	taskControlRepo    cacherepository.TaskControlRepo
	appService         appservice.AppService
	dockerManager      *docker.Manager
}

func NewAppDeploymentUC(
	db *database.DB,
	redisClient rediscache.Client,
	projectRepo repository.ProjectRepo,
	appRepo repository.AppRepo,
	deploymentRepo repository.DeploymentRepo,
	deploymentLogRepo repository.DeploymentLogRepo,
	deploymentInfoRepo cacherepository.DeploymentInfoRepo,
	taskControlRepo cacherepository.TaskControlRepo,
	appService appservice.AppService,
	dockerManager *docker.Manager,
) *AppDeploymentUC {
	return &AppDeploymentUC{
		db:                 db,
		redisClient:        redisClient,
		projectRepo:        projectRepo,
		appRepo:            appRepo,
		deploymentRepo:     deploymentRepo,
		deploymentLogRepo:  deploymentLogRepo,
		deploymentInfoRepo: deploymentInfoRepo,
		taskControlRepo:    taskControlRepo,
		appService:         appService,
		dockerManager:      dockerManager,
	}
}
