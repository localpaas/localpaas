package appuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/envvarservice"
	"github.com/localpaas/localpaas/localpaas_app/service/networkservice"
	"github.com/localpaas/localpaas/localpaas_app/service/nginxservice"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
	"github.com/localpaas/localpaas/services/docker"
)

type AppUC struct {
	db             *database.DB
	projectRepo    repository.ProjectRepo
	appRepo        repository.AppRepo
	settingRepo    repository.SettingRepo
	deploymentRepo repository.DeploymentRepo
	userService    userservice.UserService
	appService     appservice.AppService
	settingService settingservice.SettingService
	projectService projectservice.ProjectService
	networkService networkservice.NetworkService
	envVarService  envvarservice.EnvVarService
	nginxService   nginxservice.NginxService
	dockerManager  *docker.Manager
	taskQueue      taskqueue.TaskQueue
}

func NewAppUC(
	db *database.DB,
	projectRepo repository.ProjectRepo,
	appRepo repository.AppRepo,
	settingRepo repository.SettingRepo,
	deploymentRepo repository.DeploymentRepo,
	userService userservice.UserService,
	appService appservice.AppService,
	settingService settingservice.SettingService,
	projectService projectservice.ProjectService,
	networkService networkservice.NetworkService,
	envVarService envvarservice.EnvVarService,
	nginxService nginxservice.NginxService,
	dockerManager *docker.Manager,
	taskQueue taskqueue.TaskQueue,
) *AppUC {
	return &AppUC{
		db:             db,
		projectRepo:    projectRepo,
		appRepo:        appRepo,
		settingRepo:    settingRepo,
		deploymentRepo: deploymentRepo,
		userService:    userService,
		appService:     appService,
		settingService: settingService,
		projectService: projectService,
		networkService: networkService,
		envVarService:  envVarService,
		nginxService:   nginxService,
		dockerManager:  dockerManager,
		taskQueue:      taskQueue,
	}
}
