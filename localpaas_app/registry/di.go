package registry

import (
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/apphandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/clusterhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/gitsourcehandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/projecthandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/providershandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/sessionhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/settinghandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/systemhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/userhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/usersettingshandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/server"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
	"github.com/localpaas/localpaas/localpaas_app/service/envvarservice"
	"github.com/localpaas/localpaas/localpaas_app/service/lpappservice"
	"github.com/localpaas/localpaas/localpaas_app/service/networkservice"
	"github.com/localpaas/localpaas/localpaas_app/service/nginxservice"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/initializer"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/taskappdeploy"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/tasktest"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appdeploymentuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/imageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/nodeuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/gitsourceuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/discorduc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/lpappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/nginxuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/syserroruc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/taskuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc"
	"github.com/localpaas/localpaas/services/docker"
)

var Provides = []any{
	// configuration
	config.LoadConfig,

	// logger
	logging.NewZapLogger,

	// db
	database.NewDB,

	// cache
	rediscache.NewClient,

	// http server
	server.NewHTTPServer,

	// permission
	permission.NewManager,

	// Infra
	docker.New,

	// Task queue
	taskqueue.NewTaskQueue,
	initializer.NewWorkerInitializer,
	tasktest.NewExecutor,
	taskappdeploy.NewExecutor,

	// Route handler
	server.NewHandlerRegistry, // for all handler list
	handler.NewBaseHandler,
	basesettinghandler.NewBaseSettingHandler,
	authhandler.NewAuthHandler,
	clusterhandler.NewClusterHandler,
	sessionhandler.NewSessionHandler,
	userhandler.NewUserHandler,
	projecthandler.NewProjectHandler,
	apphandler.NewAppHandler,
	providershandler.NewProvidersHandler,
	settinghandler.NewSettingHandler,
	usersettingshandler.NewUserSettingsHandler,
	systemhandler.NewSystemHandler,
	gitsourcehandler.NewGitSourceHandler,

	// Use case
	syserroruc.NewSysErrorUC,
	nodeuc.NewNodeUC,
	volumeuc.NewVolumeUC,
	imageuc.NewImageUC,
	sessionuc.NewSessionUC,
	useruc.NewUserUC,
	projectuc.NewProjectUC,
	appuc.NewAppUC,
	appdeploymentuc.NewAppDeploymentUC,
	s3storageuc.NewS3StorageUC,
	sshkeyuc.NewSSHKeyUC,
	apikeyuc.NewAPIKeyUC,
	oauthuc.NewOAuthUC,
	secretuc.NewSecretUC,
	slackuc.NewSlackUC,
	discorduc.NewDiscordUC,
	registryauthuc.NewRegistryAuthUC,
	basicauthuc.NewBasicAuthUC,
	ssluc.NewSslUC,
	githubappuc.NewGithubAppUC,
	gittokenuc.NewGitTokenUC,
	nginxuc.NewNginxUC,
	lpappuc.NewLpAppUC,
	gitsourceuc.NewGitSourceUC,
	cronjobuc.NewCronJobUC,
	taskuc.NewTaskUC,

	// Service
	clusterservice.NewClusterService,
	userservice.NewUserService,
	projectservice.NewProjectService,
	appservice.NewAppService,
	networkservice.NewNetworkService,
	settingservice.NewSettingService,
	envvarservice.NewEnvVarService,
	nginxservice.NewNginxService,
	lpappservice.NewLpAppService,

	// Repo: User
	repository.NewUserRepo,
	// Repo: Permission
	repository.NewACLPermissionRepo,
	// Repo: Project
	repository.NewProjectRepo,
	repository.NewProjectTagRepo,
	repository.NewProjectSharedSettingRepo,
	// Repo: App
	repository.NewAppRepo,
	repository.NewAppTagRepo,
	// Repo: App deployment
	repository.NewDeploymentRepo,
	repository.NewDeploymentLogRepo,
	// Repo: Setting
	repository.NewSettingRepo,
	// Repo: Task
	repository.NewTaskRepo,
	// Repo: Sys error
	repository.NewSysErrorRepo,
	// Others
	repository.NewLoginTrustedDeviceRepo,

	// Cache Repo
	cacherepository.NewUserTokenRepo,
	cacherepository.NewMFAPasscodeRepo,
	cacherepository.NewTaskInfoRepo,
	cacherepository.NewTaskControlRepo,
	cacherepository.NewDeploymentInfoRepo,
}
