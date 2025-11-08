package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggoFiles "github.com/swaggo/files"
	swaggoGin "github.com/swaggo/gin-swagger"

	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/apikeyhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/apphandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/clusterhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/projecthandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/sessionhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/settingshandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/userhandler"
)

type HandlerRegistry struct {
	authHandler     *authhandler.AuthHandler
	clusterHandler  *clusterhandler.ClusterHandler
	sessionHandler  *sessionhandler.SessionHandler
	userHandler     *userhandler.UserHandler
	projectHandler  *projecthandler.ProjectHandler
	appHandler      *apphandler.AppHandler
	apiKeyHandler   *apikeyhandler.APIKeyHandler
	settingsHandler *settingshandler.SettingsHandler
}

func NewHandlerRegistry(
	authHandler *authhandler.AuthHandler,
	clusterHandler *clusterhandler.ClusterHandler,
	sessionHandler *sessionhandler.SessionHandler,
	userHandler *userhandler.UserHandler,
	projectHandler *projecthandler.ProjectHandler,
	appHandler *apphandler.AppHandler,
	apiKeyHandler *apikeyhandler.APIKeyHandler,
	settingsHandler *settingshandler.SettingsHandler,
) *HandlerRegistry {
	return &HandlerRegistry{
		authHandler:     authHandler,
		clusterHandler:  clusterHandler,
		sessionHandler:  sessionHandler,
		userHandler:     userHandler,
		projectHandler:  projectHandler,
		appHandler:      appHandler,
		apiKeyHandler:   apiKeyHandler,
		settingsHandler: settingsHandler,
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
	s.engine.Use(StaticServe("/files/user/photo", localFile(s.config.App.DataPathUserPhoto(), false)))
	// Serve the static files from the "dist-dashboard" directory at the root URL "/"
	s.engine.Use(StaticServe("/", localFile("./dist-dashboard", true)))
	// Final redirection to redirect any path to `/next=<path>` in case no matching static file found
	s.engine.Use(StaticServeRedirect("/"))

	// INTERNAL ROUTES
	basicAuthMdlw := gin.BasicAuth(gin.Accounts{
		s.config.Session.BasicAuth.Username: s.config.Session.BasicAuth.Password,
	})
	v1BasicAuth := s.engine.Group(s.config.HTTPServer.BasePath + "/internal")
	v1BasicAuth.Use(basicAuthMdlw)

	// Dev mode
	if !s.config.IsProdEnv() {
		v1BasicAuth.POST("/auth/dev-mode-login", s.handlerRegistry.sessionHandler.DevModeLogin)
	}

	// PUBLIC ROUTES
	apiV1Group := s.engine.Group(s.config.HTTPServer.BasePath)

	clusterGroup := apiV1Group.Group("/cluster")
	{ // node group
		nodeGroup := clusterGroup.Group("/nodes")
		// Nodes
		nodeGroup.GET("", s.handlerRegistry.clusterHandler.ListNode)
		nodeGroup.GET("/:nodeID", s.handlerRegistry.clusterHandler.GetNode)
		nodeGroup.DELETE("/:nodeID", s.handlerRegistry.clusterHandler.DeleteNode)
	}

	{ // session group
		sessionGroup := apiV1Group.Group("/sessions")
		// User info
		sessionGroup.GET("/me", s.handlerRegistry.sessionHandler.GetMe)
		// Session handling
		sessionGroup.POST("/refresh", s.handlerRegistry.sessionHandler.RefreshSession)
		sessionGroup.DELETE("", s.handlerRegistry.sessionHandler.DeleteSession)
		sessionGroup.POST("/delete-all", s.handlerRegistry.sessionHandler.DeleteAllSessions)
	}

	{ // auth group
		authGroup := apiV1Group.Group("/auth")
		// Login options
		authGroup.GET("/login-options", s.handlerRegistry.sessionHandler.LoginGetOptions)
		// Login with password
		authGroup.POST("/login-with-password", s.handlerRegistry.sessionHandler.LoginWithPassword)
		authGroup.POST("/login-with-passcode", s.handlerRegistry.sessionHandler.LoginWithPasscode)
		// Login with API key
		authGroup.POST("/login-with-api-key", s.handlerRegistry.sessionHandler.LoginWithAPIKey)
	}

	{ // user group
		userGroup := apiV1Group.Group("/users")
		// User info
		userGroup.GET("/base", s.handlerRegistry.userHandler.ListUserBase)
		userGroup.GET("/:userID", s.handlerRegistry.userHandler.GetUser)
		userGroup.GET("", s.handlerRegistry.userHandler.ListUser)
		// Password
		userGroup.PUT("/current/password", s.handlerRegistry.userHandler.UpdateUserPassword)
		// Profile
		userGroup.PUT("/current/profile", s.handlerRegistry.userHandler.UpdateUserProfile)
		// Status (admin API)
		userGroup.PUT("/:userID", s.handlerRegistry.userHandler.UpdateUser)
		// MFA TOTP setup
		userGroup.POST("/current/mfa/totp-begin-setup", s.handlerRegistry.userHandler.BeginMFATotpSetup)
		userGroup.POST("/current/mfa/totp-complete-setup", s.handlerRegistry.userHandler.CompleteMFATotpSetup)
		userGroup.POST("/current/mfa/totp-remove", s.handlerRegistry.userHandler.RemoveMFATotp)
		// Invite & SignUp
		userGroup.POST("/invite", s.handlerRegistry.userHandler.InviteUser)
		userGroup.POST("/signup-begin", s.handlerRegistry.userHandler.BeginUserSignup)
		userGroup.POST("/signup-complete", s.handlerRegistry.userHandler.CompleteUserSignup)

		// API key group
		apiKeyGroup := userGroup.Group("/current/api-keys")
		{
			// Info
			apiKeyGroup.GET("/:ID", s.handlerRegistry.apiKeyHandler.GetAPIKey)
			apiKeyGroup.GET("", s.handlerRegistry.apiKeyHandler.ListAPIKey)
			// Creation & Update
			apiKeyGroup.POST("", s.handlerRegistry.apiKeyHandler.CreateAPIKey)
			apiKeyGroup.DELETE("/:ID", s.handlerRegistry.apiKeyHandler.DeleteAPIKey)
		}
	}

	projectGroup := apiV1Group.Group("/projects")
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
		// Settings
		projectGroup.GET("/:projectID/settings", s.handlerRegistry.projectHandler.GetProjectSettings)
		projectGroup.PUT("/:projectID/settings", s.handlerRegistry.projectHandler.UpdateProjectSettings)
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
		appGroup.GET("/:appID/settings", s.handlerRegistry.appHandler.GetAppSettings)
		appGroup.PUT("/:appID/settings", s.handlerRegistry.appHandler.UpdateAppSettings)
	}

	settingGroup := apiV1Group.Group("/settings")

	{ // ssh key group
		oauthGroup := settingGroup.Group("/oauth")
		// Info
		oauthGroup.GET("/:ID", s.handlerRegistry.settingsHandler.GetOAuth)
		oauthGroup.GET("", s.handlerRegistry.settingsHandler.ListOAuth)
		// Creation & Update
		oauthGroup.POST("", s.handlerRegistry.settingsHandler.CreateOAuth)
		oauthGroup.PUT("/:ID", s.handlerRegistry.settingsHandler.UpdateOAuth)
		oauthGroup.DELETE("/:ID", s.handlerRegistry.settingsHandler.DeleteOAuth)
	}

	{ // s3 storage group
		s3StorageGroup := settingGroup.Group("/s3-storages")
		// Info
		s3StorageGroup.GET("/:ID", s.handlerRegistry.settingsHandler.GetS3Storage)
		s3StorageGroup.GET("", s.handlerRegistry.settingsHandler.ListS3Storage)
		// Creation & Update
		s3StorageGroup.POST("", s.handlerRegistry.settingsHandler.CreateS3Storage)
		s3StorageGroup.PUT("/:ID", s.handlerRegistry.settingsHandler.UpdateS3Storage)
		s3StorageGroup.DELETE("/:ID", s.handlerRegistry.settingsHandler.DeleteS3Storage)
	}

	{ // ssh key group
		sshKeyGroup := settingGroup.Group("/ssh-keys")
		// Info
		sshKeyGroup.GET("/:ID", s.handlerRegistry.settingsHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.settingsHandler.ListSSHKey)
		// Creation & Update
		sshKeyGroup.POST("", s.handlerRegistry.settingsHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:ID", s.handlerRegistry.settingsHandler.UpdateSSHKey)
		sshKeyGroup.DELETE("/:ID", s.handlerRegistry.settingsHandler.DeleteSSHKey)
	}
}

func routePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func routeNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, "not found")
}
