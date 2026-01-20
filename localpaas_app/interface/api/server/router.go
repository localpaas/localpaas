package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggoFiles "github.com/swaggo/files"
	swaggoGin "github.com/swaggo/gin-swagger"

	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/apphandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/clusterhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/gitsourcehandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/projecthandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/providershandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/sessionhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/systemhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/userhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/usersettingshandler"
)

type HandlerRegistry struct {
	authHandler         *authhandler.AuthHandler
	clusterHandler      *clusterhandler.ClusterHandler
	sessionHandler      *sessionhandler.SessionHandler
	userHandler         *userhandler.UserHandler
	projectHandler      *projecthandler.ProjectHandler
	appHandler          *apphandler.AppHandler
	providersHandler    *providershandler.ProvidersHandler
	userSettingsHandler *usersettingshandler.UserSettingsHandler
	systemHandler       *systemhandler.SystemHandler
	gitSourceHandler    *gitsourcehandler.GitSourceHandler
}

func NewHandlerRegistry(
	authHandler *authhandler.AuthHandler,
	clusterHandler *clusterhandler.ClusterHandler,
	sessionHandler *sessionhandler.SessionHandler,
	userHandler *userhandler.UserHandler,
	projectHandler *projecthandler.ProjectHandler,
	appHandler *apphandler.AppHandler,
	providersHandler *providershandler.ProvidersHandler,
	userSettingsHandler *usersettingshandler.UserSettingsHandler,
	systemHandler *systemhandler.SystemHandler,
	gitSourceHandler *gitsourcehandler.GitSourceHandler,
) *HandlerRegistry {
	return &HandlerRegistry{
		authHandler:         authHandler,
		clusterHandler:      clusterHandler,
		sessionHandler:      sessionHandler,
		userHandler:         userHandler,
		projectHandler:      projectHandler,
		appHandler:          appHandler,
		providersHandler:    providersHandler,
		userSettingsHandler: userSettingsHandler,
		systemHandler:       systemHandler,
		gitSourceHandler:    gitSourceHandler,
	}
}

//nolint:funlen
func (s *HTTPServer) registerRoutes() {
	s.engine.GET("/_/ping", routePing)
	s.engine.NoRoute(routeNotFound)

	// Swagger server
	if !s.config.IsProdEnv() {
		s.engine.Use(StaticServe("/docs", localFile("./docs", false)))
		s.engine.GET("/swagger/*any", swaggoGin.WrapHandler(swaggoFiles.Handler,
			swaggoGin.URL("/docs/openapi/swagger.json")))
	}

	// STATIC FILES
	s.engine.Use(StaticServe(s.config.HttpPathUserPhoto(), localFile(s.config.DataPathUserPhoto(), false)))
	// Serve the static files from the "dist-dashboard" directory at the root URL "/"
	s.engine.Use(StaticServe("/", localFile("./dist-dashboard", true)))
	// Final redirection to redirect any path to `/next=<path>` in case no matching static file found
	s.engine.Use(StaticServeRedirect("/"))

	// INTERNAL ROUTES
	basicAuthMdlw := gin.BasicAuth(gin.Accounts{
		s.config.Session.BasicAuthUsername: s.config.Session.BasicAuthPassword,
	})
	v1BasicAuth := s.engine.Group(s.config.HTTPServer.BasePath + "/internal")
	v1BasicAuth.Use(basicAuthMdlw)

	// Dev mode
	if !s.config.IsProdEnv() {
		v1BasicAuth.POST("/auth/dev-mode-login", s.handlerRegistry.sessionHandler.DevModeLogin)
	}

	// PUBLIC ROUTES
	apiGroup := s.engine.Group(s.config.HTTPServer.BasePath)

	{ // session group
		sessionGroup := apiGroup.Group("/sessions")
		// User info
		sessionGroup.GET("/me", s.handlerRegistry.sessionHandler.GetMe)
		// Session handling
		sessionGroup.POST("/refresh", s.handlerRegistry.sessionHandler.RefreshSession)
		sessionGroup.DELETE("", s.handlerRegistry.sessionHandler.DeleteSession)
		sessionGroup.POST("/delete-all", s.handlerRegistry.sessionHandler.DeleteAllSessions)
	}

	{ // auth group
		authGroup := apiGroup.Group("/auth")
		// Login options
		authGroup.GET("/login-options", s.handlerRegistry.sessionHandler.LoginGetOptions)
		// Login with password
		authGroup.POST("/login-with-password", s.handlerRegistry.sessionHandler.LoginWithPassword)
		authGroup.POST("/login-with-passcode", s.handlerRegistry.sessionHandler.LoginWithPasscode)
		// Login with API key
		authGroup.POST("/login-with-api-key", s.handlerRegistry.sessionHandler.LoginWithAPIKey)
		// Login via SSO
		authGroup.GET("/sso/:provider", s.handlerRegistry.sessionHandler.SSOOAuthBegin)
		authGroup.GET("/sso/callback/:provider", s.handlerRegistry.sessionHandler.SSOOAuthCallback)
		authGroup.POST("/sso/callback/:provider", s.handlerRegistry.sessionHandler.SSOOAuthCallback)
	}

	userGroup := apiGroup.Group("/users")
	{ // user group
		// User info
		userGroup.GET("/base", s.handlerRegistry.userHandler.ListUserBase)
		userGroup.GET("/:userID", s.handlerRegistry.userHandler.GetUser)
		userGroup.GET("", s.handlerRegistry.userHandler.ListUser)
		// Password
		userGroup.PUT("/current/password", s.handlerRegistry.userHandler.UpdateUserPassword)
		userGroup.POST("/:userID/password/request-reset", s.handlerRegistry.userHandler.RequestResetPassword)
		userGroup.POST("/:userID/password/reset", s.handlerRegistry.userHandler.ResetPassword)
		// Profile
		userGroup.PUT("/current/profile", s.handlerRegistry.userHandler.UpdateUserProfile)
		// Update (admin API)
		userGroup.PUT("/:userID", s.handlerRegistry.userHandler.UpdateUser)
		userGroup.DELETE("/:userID", s.handlerRegistry.userHandler.DeleteUser)
		// MFA TOTP setup
		userGroup.POST("/current/mfa/totp-begin-setup", s.handlerRegistry.userHandler.BeginMFATotpSetup)
		userGroup.POST("/current/mfa/totp-complete-setup", s.handlerRegistry.userHandler.CompleteMFATotpSetup)
		userGroup.POST("/current/mfa/totp-remove", s.handlerRegistry.userHandler.RemoveMFATotp)
		// Invite & SignUp
		userGroup.POST("/invite", s.handlerRegistry.userHandler.InviteUser)
		userGroup.POST("/signup-begin", s.handlerRegistry.userHandler.BeginUserSignup)
		userGroup.POST("/signup-complete", s.handlerRegistry.userHandler.CompleteUserSignup)
	}

	// User settings group
	userSettingGroup := userGroup.Group("/current/settings")

	{ // API key group
		apiKeyGroup := userSettingGroup.Group("/api-keys")
		// Info
		apiKeyGroup.GET("/:id", s.handlerRegistry.userSettingsHandler.GetAPIKey)
		apiKeyGroup.GET("", s.handlerRegistry.userSettingsHandler.ListAPIKey)
		// Creation & Update
		apiKeyGroup.POST("", s.handlerRegistry.userSettingsHandler.CreateAPIKey)
		apiKeyGroup.PUT("/:id/meta", s.handlerRegistry.userSettingsHandler.UpdateAPIKeyMeta)
		apiKeyGroup.DELETE("/:id", s.handlerRegistry.userSettingsHandler.DeleteAPIKey)
	}

	clusterGroup := apiGroup.Group("/cluster")
	{ // node group
		nodeGroup := clusterGroup.Group("/nodes")
		// Nodes
		nodeGroup.GET("", s.handlerRegistry.clusterHandler.ListNode)
		nodeGroup.GET("/:nodeID", s.handlerRegistry.clusterHandler.GetNode)
		nodeGroup.GET("/:nodeID/inspect", s.handlerRegistry.clusterHandler.GetNodeInspection)
		nodeGroup.PUT("/:nodeID", s.handlerRegistry.clusterHandler.UpdateNode)
		nodeGroup.DELETE("/:nodeID", s.handlerRegistry.clusterHandler.DeleteNode)
		// Node join
		nodeGroup.POST("/join", s.handlerRegistry.clusterHandler.JoinNode)
		nodeGroup.GET("/join-command", s.handlerRegistry.clusterHandler.GetNodeJoinCommand)
	}
	{ // volume group
		volumeGroup := clusterGroup.Group("/volumes")
		// Volumes
		volumeGroup.GET("", s.handlerRegistry.clusterHandler.ListVolume)
		volumeGroup.GET("/:volumeID", s.handlerRegistry.clusterHandler.GetVolume)
		volumeGroup.GET("/:volumeID/inspect", s.handlerRegistry.clusterHandler.GetVolumeInspection)
		volumeGroup.POST("", s.handlerRegistry.clusterHandler.CreateVolume)
		volumeGroup.DELETE("/:volumeID", s.handlerRegistry.clusterHandler.DeleteVolume)
	}
	{ // image group
		imageGroup := clusterGroup.Group("/images")
		// Volumes
		imageGroup.GET("", s.handlerRegistry.clusterHandler.ListImage)
		imageGroup.GET("/:imageID", s.handlerRegistry.clusterHandler.GetImage)
		imageGroup.GET("/:imageID/inspect", s.handlerRegistry.clusterHandler.GetImageInspection)
		imageGroup.POST("", s.handlerRegistry.clusterHandler.CreateImage)
		imageGroup.DELETE("/:imageID", s.handlerRegistry.clusterHandler.DeleteImage)
	}

	projectGroup := apiGroup.Group("/projects")
	{ // project group
		// Project info
		projectGroup.GET("/base", s.handlerRegistry.projectHandler.ListProjectBase)
		projectGroup.GET("/:projectID", s.handlerRegistry.projectHandler.GetProject)
		projectGroup.GET("", s.handlerRegistry.projectHandler.ListProject)
		// Creation & Update
		projectGroup.POST("", s.handlerRegistry.projectHandler.CreateProject)
		projectGroup.DELETE("/:projectID", s.handlerRegistry.projectHandler.DeleteProject)
		// Tags
		projectGroup.POST("/:projectID/tags", s.handlerRegistry.projectHandler.CreateProjectTag)
		projectGroup.POST("/:projectID/tags/delete", s.handlerRegistry.projectHandler.DeleteProjectTags)
		// Env vars
		projectGroup.GET("/:projectID/env-vars", s.handlerRegistry.projectHandler.GetProjectEnvVars)
		projectGroup.PUT("/:projectID/env-vars", s.handlerRegistry.projectHandler.UpdateProjectEnvVars)
		// Secrets
		projectGroup.GET("/:projectID/secrets", s.handlerRegistry.projectHandler.ListProjectSecrets)
		projectGroup.POST("/:projectID/secrets", s.handlerRegistry.projectHandler.CreateProjectSecret)
		projectGroup.DELETE("/:projectID/secrets/:id", s.handlerRegistry.projectHandler.DeleteProjectSecret)
		// Cron jobs
		projectGroup.GET("/:projectID/cron-jobs", s.handlerRegistry.projectHandler.ListCronJob)
		projectGroup.GET("/:projectID/cron-jobs/:id", s.handlerRegistry.projectHandler.GetCronJob)
		projectGroup.POST("/:projectID/cron-jobs", s.handlerRegistry.projectHandler.CreateCronJob)
		projectGroup.PUT("/:projectID/cron-jobs/:id", s.handlerRegistry.projectHandler.UpdateCronJob)
		projectGroup.PUT("/:projectID/cron-jobs/:id/meta", s.handlerRegistry.projectHandler.UpdateCronJobMeta)
		projectGroup.DELETE("/:projectID/cron-jobs/:id", s.handlerRegistry.projectHandler.DeleteCronJob)
	}

	projectProviderGroup := projectGroup.Group("/:projectID/providers")

	{ // github-app group
		githubAppGroup := projectProviderGroup.Group("/github-apps")
		// Info
		githubAppGroup.GET("/:id", s.handlerRegistry.projectHandler.GetGithubApp)
		githubAppGroup.GET("", s.handlerRegistry.projectHandler.ListGithubApp)
		// Creation & Update
		githubAppGroup.POST("", s.handlerRegistry.projectHandler.CreateGithubApp)
		githubAppGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateGithubApp)
		githubAppGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateGithubAppMeta)
		githubAppGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteGithubApp)
	}

	{ // git-token group
		gitTokenGroup := projectProviderGroup.Group("/git-tokens")
		// Info
		gitTokenGroup.GET("/:id", s.handlerRegistry.projectHandler.GetGitToken)
		gitTokenGroup.GET("", s.handlerRegistry.projectHandler.ListGitToken)
		// Creation & Update
		gitTokenGroup.POST("", s.handlerRegistry.projectHandler.CreateGitToken)
		gitTokenGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateGitToken)
		gitTokenGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateGitTokenMeta)
		gitTokenGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteGitToken)
	}

	{ // s3 storage group
		s3StorageGroup := projectProviderGroup.Group("/s3-storages")
		// Info
		s3StorageGroup.GET("/:id", s.handlerRegistry.projectHandler.GetS3Storage)
		s3StorageGroup.GET("", s.handlerRegistry.projectHandler.ListS3Storage)
		// Creation & Update
		s3StorageGroup.POST("", s.handlerRegistry.projectHandler.CreateS3Storage)
		s3StorageGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateS3Storage)
		s3StorageGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateS3StorageMeta)
		s3StorageGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteS3Storage)
	}

	{ // ssh key group
		sshKeyGroup := projectProviderGroup.Group("/ssh-keys")
		// Info
		sshKeyGroup.GET("/:id", s.handlerRegistry.projectHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.projectHandler.ListSSHKey)
		// Creation & Update
		sshKeyGroup.POST("", s.handlerRegistry.projectHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateSSHKey)
		sshKeyGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateSSHKeyMeta)
		sshKeyGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteSSHKey)
	}

	{ // slack group
		slackGroup := projectProviderGroup.Group("/slack")
		// Info
		slackGroup.GET("/:id", s.handlerRegistry.projectHandler.GetSlack)
		slackGroup.GET("", s.handlerRegistry.projectHandler.ListSlack)
		// Creation & Update
		slackGroup.POST("", s.handlerRegistry.projectHandler.CreateSlack)
		slackGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateSlack)
		slackGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateSlackMeta)
		slackGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteSlack)
	}

	{ // discord group
		discordGroup := projectProviderGroup.Group("/discord")
		// Info
		discordGroup.GET("/:id", s.handlerRegistry.projectHandler.GetDiscord)
		discordGroup.GET("", s.handlerRegistry.projectHandler.ListDiscord)
		// Creation & Update
		discordGroup.POST("", s.handlerRegistry.projectHandler.CreateDiscord)
		discordGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateDiscord)
		discordGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateDiscordMeta)
		discordGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteDiscord)
	}

	{ // registry auth group
		registryAuthGroup := projectProviderGroup.Group("/registry-auth")
		// Info
		registryAuthGroup.GET("/:id", s.handlerRegistry.projectHandler.GetRegistryAuth)
		registryAuthGroup.GET("", s.handlerRegistry.projectHandler.ListRegistryAuth)
		// Creation & Update
		registryAuthGroup.POST("", s.handlerRegistry.projectHandler.CreateRegistryAuth)
		registryAuthGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateRegistryAuth)
		registryAuthGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateRegistryAuthMeta)
		registryAuthGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteRegistryAuth)
	}

	{ // basic auth group
		basicAuthGroup := projectProviderGroup.Group("/basic-auth")
		// Info
		basicAuthGroup.GET("/:id", s.handlerRegistry.projectHandler.GetBasicAuth)
		basicAuthGroup.GET("", s.handlerRegistry.projectHandler.ListBasicAuth)
		// Creation & Update
		basicAuthGroup.POST("", s.handlerRegistry.projectHandler.CreateBasicAuth)
		basicAuthGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateBasicAuth)
		basicAuthGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateBasicAuthMeta)
		basicAuthGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteBasicAuth)
	}

	{ // ssl group
		sslGroup := projectProviderGroup.Group("/ssls")
		// Info
		sslGroup.GET("/:id", s.handlerRegistry.projectHandler.GetSsl)
		sslGroup.GET("", s.handlerRegistry.projectHandler.ListSsl)
		// Creation & Update
		sslGroup.POST("", s.handlerRegistry.projectHandler.CreateSsl)
		sslGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateSsl)
		sslGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateSslMeta)
		sslGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteSsl)
	}

	appGroup := projectGroup.Group("/:projectID/apps")
	{ // app group
		// App info
		appGroup.GET("/base", s.handlerRegistry.appHandler.ListAppBase)
		appGroup.GET("/:appID", s.handlerRegistry.appHandler.GetApp)
		appGroup.GET("", s.handlerRegistry.appHandler.ListApp)
		// Creation & Update
		appGroup.POST("", s.handlerRegistry.appHandler.CreateApp)
		appGroup.DELETE("/:appID", s.handlerRegistry.appHandler.DeleteApp)
		// Tags
		appGroup.POST("/:appID/tags", s.handlerRegistry.appHandler.CreateAppTag)
		appGroup.POST("/:appID/tags/delete", s.handlerRegistry.appHandler.DeleteAppTags)
		// Settings
		appGroup.GET("/:appID/service-spec", s.handlerRegistry.appHandler.GetAppServiceSpec)
		appGroup.PUT("/:appID/service-spec", s.handlerRegistry.appHandler.UpdateAppServiceSpec)
		appGroup.GET("/:appID/deployment-settings", s.handlerRegistry.appHandler.GetAppDeploymentSettings)
		appGroup.PUT("/:appID/deployment-settings", s.handlerRegistry.appHandler.UpdateAppDeploymentSettings)
		appGroup.GET("/:appID/http-settings", s.handlerRegistry.appHandler.GetAppHttpSettings)
		appGroup.PUT("/:appID/http-settings", s.handlerRegistry.appHandler.UpdateAppHttpSettings)
		// Env vars
		appGroup.GET("/:appID/env-vars", s.handlerRegistry.appHandler.GetAppEnvVars)
		appGroup.PUT("/:appID/env-vars", s.handlerRegistry.appHandler.UpdateAppEnvVars)
		// Secrets
		appGroup.GET("/:appID/secrets", s.handlerRegistry.appHandler.ListAppSecrets)
		appGroup.POST("/:appID/secrets", s.handlerRegistry.appHandler.CreateAppSecret)
		appGroup.DELETE("/:appID/secrets/:secretID", s.handlerRegistry.appHandler.DeleteAppSecret)
		// Domain SSL
		appGroup.POST("/:appID/ssl/obtain", s.handlerRegistry.appHandler.ObtainDomainSsl)
		// Logs
		appGroup.GET("/:appID/runtime-logs", func(ctx *gin.Context) {
			s.handlerRegistry.appHandler.GetAppRuntimeLogs(ctx, s.websocket)
		})
	}

	appDeploymentGroup := appGroup.Group("/:appID/deployments")
	{ // app deployment group
		// Info
		appDeploymentGroup.GET("/:deploymentID", s.handlerRegistry.appHandler.GetAppDeployment)
		appDeploymentGroup.GET("", s.handlerRegistry.appHandler.ListAppDeployment)
		// Cancel
		appDeploymentGroup.POST("/:deploymentID/cancel", s.handlerRegistry.appHandler.CancelAppDeployment)
		// Logs
		appDeploymentGroup.GET("/:deploymentID/logs", func(ctx *gin.Context) {
			s.handlerRegistry.appHandler.GetAppDeploymentLogs(ctx, s.websocket)
		})
	}

	providerGroup := apiGroup.Group("/providers")

	{ // oauth group
		oauthGroup := providerGroup.Group("/oauth")
		// Info
		oauthGroup.GET("/:id", s.handlerRegistry.providersHandler.GetOAuth)
		oauthGroup.GET("", s.handlerRegistry.providersHandler.ListOAuth)
		// Creation & Update
		oauthGroup.POST("", s.handlerRegistry.providersHandler.CreateOAuth)
		oauthGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateOAuth)
		oauthGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateOAuthMeta)
		oauthGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteOAuth)
	}

	{ // github-app group
		githubAppGroup := providerGroup.Group("/github-apps")
		// Info
		githubAppGroup.GET("/:id", s.handlerRegistry.providersHandler.GetGithubApp)
		githubAppGroup.GET("", s.handlerRegistry.providersHandler.ListGithubApp)
		// Creation & Update
		githubAppGroup.POST("", s.handlerRegistry.providersHandler.CreateGithubApp)
		githubAppGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateGithubApp)
		githubAppGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateGithubAppMeta)
		githubAppGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteGithubApp)
		// Test connection
		githubAppGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestGithubAppConn)
		// Installation
		githubAppGroup.POST("/installations/list", s.handlerRegistry.providersHandler.ListAppInstallation)
	}

	{ // git-token group
		gitTokenGroup := providerGroup.Group("/git-tokens")
		// Info
		gitTokenGroup.GET("/:id", s.handlerRegistry.providersHandler.GetGitToken)
		gitTokenGroup.GET("", s.handlerRegistry.providersHandler.ListGitToken)
		// Creation & Update
		gitTokenGroup.POST("", s.handlerRegistry.providersHandler.CreateGitToken)
		gitTokenGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateGitToken)
		gitTokenGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateGitTokenMeta)
		gitTokenGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteGitToken)
		// Test connection
		gitTokenGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestGitTokenConn)
	}

	{ // s3 storage group
		s3StorageGroup := providerGroup.Group("/s3-storages")
		// Info
		s3StorageGroup.GET("/:id", s.handlerRegistry.providersHandler.GetS3Storage)
		s3StorageGroup.GET("", s.handlerRegistry.providersHandler.ListS3Storage)
		// Creation & Update
		s3StorageGroup.POST("", s.handlerRegistry.providersHandler.CreateS3Storage)
		s3StorageGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateS3Storage)
		s3StorageGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateS3StorageMeta)
		s3StorageGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteS3Storage)
		// Test connection
		s3StorageGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestS3StorageConn)
	}

	{ // ssh key group
		sshKeyGroup := providerGroup.Group("/ssh-keys")
		// Info
		sshKeyGroup.GET("/:id", s.handlerRegistry.providersHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.providersHandler.ListSSHKey)
		// Creation & Update
		sshKeyGroup.POST("", s.handlerRegistry.providersHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateSSHKey)
		sshKeyGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateSSHKeyMeta)
		sshKeyGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteSSHKey)
	}

	{ // secrets group
		secretGroup := providerGroup.Group("/secrets")
		// Info
		secretGroup.GET("", s.handlerRegistry.providersHandler.ListSecret)
		// Creation & Update
		secretGroup.POST("", s.handlerRegistry.providersHandler.CreateSecret)
		secretGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateSecretMeta)
		secretGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteSecret)
	}

	{ // slack group
		slackGroup := providerGroup.Group("/slack")
		// Info
		slackGroup.GET("/:id", s.handlerRegistry.providersHandler.GetSlack)
		slackGroup.GET("", s.handlerRegistry.providersHandler.ListSlack)
		// Creation & Update
		slackGroup.POST("", s.handlerRegistry.providersHandler.CreateSlack)
		slackGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateSlack)
		slackGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateSlackMeta)
		slackGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteSlack)
		// Test connection
		slackGroup.POST("/test-send-msg", s.handlerRegistry.providersHandler.TestSendSlackMsg)
	}

	{ // discord group
		discordGroup := providerGroup.Group("/discord")
		// Info
		discordGroup.GET("/:id", s.handlerRegistry.providersHandler.GetDiscord)
		discordGroup.GET("", s.handlerRegistry.providersHandler.ListDiscord)
		// Creation & Update
		discordGroup.POST("", s.handlerRegistry.providersHandler.CreateDiscord)
		discordGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateDiscord)
		discordGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateDiscordMeta)
		discordGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteDiscord)
		// Test connection
		discordGroup.POST("/test-send-msg", s.handlerRegistry.providersHandler.TestSendDiscordMsg)
	}

	{ // registry auth group
		registryAuthGroup := providerGroup.Group("/registry-auth")
		// Info
		registryAuthGroup.GET("/:id", s.handlerRegistry.providersHandler.GetRegistryAuth)
		registryAuthGroup.GET("", s.handlerRegistry.providersHandler.ListRegistryAuth)
		// Creation & Update
		registryAuthGroup.POST("", s.handlerRegistry.providersHandler.CreateRegistryAuth)
		registryAuthGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateRegistryAuth)
		registryAuthGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateRegistryAuthMeta)
		registryAuthGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteRegistryAuth)
		// Test connection
		registryAuthGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestRegistryAuthConn)
	}

	{ // basic auth group
		basicAuthGroup := providerGroup.Group("/basic-auth")
		// Info
		basicAuthGroup.GET("/:id", s.handlerRegistry.providersHandler.GetBasicAuth)
		basicAuthGroup.GET("", s.handlerRegistry.providersHandler.ListBasicAuth)
		// Creation & Update
		basicAuthGroup.POST("", s.handlerRegistry.providersHandler.CreateBasicAuth)
		basicAuthGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateBasicAuth)
		basicAuthGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateBasicAuthMeta)
		basicAuthGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteBasicAuth)
	}

	{ // ssl group
		sslGroup := providerGroup.Group("/ssls")
		// Info
		sslGroup.GET("/:id", s.handlerRegistry.providersHandler.GetSsl)
		sslGroup.GET("", s.handlerRegistry.providersHandler.ListSsl)
		// Creation & Update
		sslGroup.POST("", s.handlerRegistry.providersHandler.CreateSsl)
		sslGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateSsl)
		sslGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateSslMeta)
		sslGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteSsl)
	}

	{ // cron-job group
		cronJobGroup := providerGroup.Group("/cron-jobs")
		// Info
		cronJobGroup.GET("/:id", s.handlerRegistry.providersHandler.GetCronJob)
		cronJobGroup.GET("", s.handlerRegistry.providersHandler.ListCronJob)
		// Creation & Update
		cronJobGroup.POST("", s.handlerRegistry.providersHandler.CreateCronJob)
		cronJobGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateCronJob)
		cronJobGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateCronJobMeta)
		cronJobGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteCronJob)
	}

	systemGroup := apiGroup.Group("/system")

	{ // task group
		taskGroup := systemGroup.Group("/tasks")
		taskGroup.GET("", s.handlerRegistry.systemHandler.ListTask)
		taskGroup.GET("/:id", s.handlerRegistry.systemHandler.GetTask)
		taskGroup.PUT("/:id/meta", s.handlerRegistry.systemHandler.UpdateTaskMeta)
		taskGroup.POST("/:id/cancel", s.handlerRegistry.systemHandler.CancelTask)
	}

	{ // error group
		errorGroup := systemGroup.Group("/errors")
		errorGroup.GET("", s.handlerRegistry.systemHandler.ListSysError)
		errorGroup.GET("/:id", s.handlerRegistry.systemHandler.GetSysError)
		errorGroup.DELETE("/:id", s.handlerRegistry.systemHandler.DeleteSysError)
	}

	{ // nginx group
		nginxGroup := systemGroup.Group("/nginx")
		// Process
		nginxGroup.POST("/restart", s.handlerRegistry.systemHandler.RestartNginx)
		// Config
		nginxGroup.POST("/config/reload", s.handlerRegistry.systemHandler.ReloadNginxConfig)
		nginxGroup.POST("/config/reset", s.handlerRegistry.systemHandler.ResetNginxConfig)
	}

	{ // localpaas app group
		lpAppGroup := systemGroup.Group("/localpaas")
		// Process
		lpAppGroup.POST("/restart", s.handlerRegistry.systemHandler.RestartLocalPaasApp)
		// Config
		lpAppGroup.POST("/config/reload", s.handlerRegistry.systemHandler.ReloadLocalPaasAppConfig)
	}

	{ // git source group
		gitSourceGroup := apiGroup.Group("/git-source")
		// Repo
		gitSourceGroup.GET("/:settingID/repositories", s.handlerRegistry.gitSourceHandler.ListGitRepo)
	}
}

func routePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func routeNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, "not found")
}
