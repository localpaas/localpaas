package appservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/nginxservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/services/docker"
)

type AppService interface {
	LoadApp(ctx context.Context, db database.IDB, projectID, appID string,
		requireProjectActive, requireAppActive bool, extraOpts ...bunex.SelectQueryOption) (
		*entity.App, error)
	LoadAppByToken(ctx context.Context, db database.IDB, appToken string,
		requireProjectActive, requireAppActive bool, extraOpts ...bunex.SelectQueryOption) (
		*entity.App, error)

	PersistAppData(ctx context.Context, db database.IDB, data *PersistingAppData) error
	DeleteApp(ctx context.Context, app *entity.App) error

	ServiceInspect(ctx context.Context, serviceID string, caching bool) (*swarm.Service, error)
	ServiceUpdate(ctx context.Context, serviceID string, version *swarm.Version, service *swarm.ServiceSpec,
		options ...docker.ServiceUpdateOption) (*swarm.ServiceUpdateResponse, error)

	LoadSettings(ctx context.Context, db database.IDB, app *entity.App, settingIDs []string,
		requireActive bool) (map[string]*entity.Setting, error)
	LoadReferenceSettings(ctx context.Context, db database.IDB, app *entity.App, requireActive bool,
		appSettings ...*entity.Setting) (map[string]*entity.Setting, error)

	EnsureSSLConfigFiles(sslIDs []string, forceRecreate bool,
		refSettingMap map[string]*entity.Setting) error
	EnsureBasicAuthConfigFiles(basicAuthIDs []string, forceRecreate bool,
		refSettingMap map[string]*entity.Setting) error

	CreateDeployment(app *entity.App, deploymentSettings *entity.AppDeploymentSettings) (
		*entity.Deployment, *entity.Task, error)
}

func NewAppService(
	appRepo repository.AppRepo,
	appTagRepo repository.AppTagRepo,
	settingRepo repository.SettingRepo,
	deploymentRepo repository.DeploymentRepo,
	taskRepo repository.TaskRepo,
	deploymentInfoRepo cacherepository.DeploymentInfoRepo,
	permissionManager permission.Manager,
	userService userservice.UserService,
	nginxService nginxservice.NginxService,
	dockerManager *docker.Manager,
) AppService {
	return &appService{
		appRepo:            appRepo,
		appTagRepo:         appTagRepo,
		settingRepo:        settingRepo,
		deploymentRepo:     deploymentRepo,
		taskRepo:           taskRepo,
		deploymentInfoRepo: deploymentInfoRepo,
		permissionManager:  permissionManager,
		userService:        userService,
		nginxService:       nginxService,
		dockerManager:      dockerManager,
	}
}

type appService struct {
	appRepo            repository.AppRepo
	appTagRepo         repository.AppTagRepo
	settingRepo        repository.SettingRepo
	deploymentRepo     repository.DeploymentRepo
	taskRepo           repository.TaskRepo
	deploymentInfoRepo cacherepository.DeploymentInfoRepo
	permissionManager  permission.Manager
	userService        userservice.UserService
	nginxService       nginxservice.NginxService
	dockerManager      *docker.Manager
}
