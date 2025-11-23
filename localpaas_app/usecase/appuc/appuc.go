package appuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/envvarservice"
	"github.com/localpaas/localpaas/localpaas_app/service/nginxservice"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
	"github.com/localpaas/localpaas/services/letsencrypt"
)

type AppUC struct {
	db                *database.DB
	userRepo          repository.UserRepo
	projectRepo       repository.ProjectRepo
	appRepo           repository.AppRepo
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
	userService       userservice.UserService
	appService        appservice.AppService
	projectService    projectservice.ProjectService
	envVarService     envvarservice.EnvVarService
	nginxService      nginxservice.NginxService
	dockerManager     *docker.Manager
	letsencryptClient *letsencrypt.Client
}

func NewAppUC(
	db *database.DB,
	userRepo repository.UserRepo,
	projectRepo repository.ProjectRepo,
	appRepo repository.AppRepo,
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
	appService appservice.AppService,
	projectService projectservice.ProjectService,
	envVarService envvarservice.EnvVarService,
	nginxService nginxservice.NginxService,
	dockerManager *docker.Manager,
	letsencryptClient *letsencrypt.Client,
) *AppUC {
	return &AppUC{
		db:                db,
		userRepo:          userRepo,
		projectRepo:       projectRepo,
		appRepo:           appRepo,
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
		userService:       userService,
		appService:        appService,
		projectService:    projectService,
		envVarService:     envVarService,
		nginxService:      nginxService,
		dockerManager:     dockerManager,
		letsencryptClient: letsencryptClient,
	}
}
