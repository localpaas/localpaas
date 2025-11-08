package registry

import (
	"github.com/localpaas/localpaas/infrastructure/docker"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/apikeyhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/apphandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/clusterhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/projecthandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/s3storagehandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/sessionhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/settingshandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/sshkeyhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/userhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/server"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apikeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/s3storageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sshkeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc"
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

	// Route handler
	server.NewHandlerRegistry, // for all handler list
	authhandler.NewAuthHandler,
	clusterhandler.NewClusterHandler,
	sessionhandler.NewSessionHandler,
	userhandler.NewUserHandler,
	projecthandler.NewProjectHandler,
	apphandler.NewAppHandler,
	s3storagehandler.NewS3StorageHandler,
	sshkeyhandler.NewSSHKeyHandler,
	apikeyhandler.NewAPIKeyHandler,
	settingshandler.NewSettingsHandler,

	// Use case
	clusteruc.NewClusterUC,
	sessionuc.NewSessionUC,
	useruc.NewUserUC,
	projectuc.NewProjectUC,
	appuc.NewAppUC,
	s3storageuc.NewS3StorageUC,
	sshkeyuc.NewSSHKeyUC,
	apikeyuc.NewAPIKeyUC,
	oauthuc.NewOAuthUC,

	// Service
	clusterservice.NewClusterService,
	userservice.NewUserService,
	projectservice.NewProjectService,
	appservice.NewAppService,
	settingservice.NewSettingService,

	// Repo: User
	repository.NewUserRepo,
	// Repo: Project
	repository.NewProjectRepo,
	repository.NewProjectTagRepo,
	// Repo: App
	repository.NewAppRepo,
	repository.NewAppTagRepo,
	// Repo: Permission
	repository.NewACLPermissionRepo,
	// Repo: Setting
	repository.NewSettingRepo,
	// Others
	repository.NewLoginTrustedDeviceRepo,

	// Cache Repo
	cacherepository.NewUserTokenRepo,
	cacherepository.NewMFAPasscodeRepo,
}
