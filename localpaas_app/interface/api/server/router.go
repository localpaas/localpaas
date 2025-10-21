package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggoFiles "github.com/swaggo/files"
	swaggoGin "github.com/swaggo/gin-swagger"

	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/apphandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/projecthandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/sessionhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/userhandler"
)

type HandlerRegistry struct {
	authHandler    *authhandler.AuthHandler
	sessionHandler *sessionhandler.SessionHandler
	userHandler    *userhandler.UserHandler
	projectHandler *projecthandler.ProjectHandler
	appHandler     *apphandler.AppHandler
}

func NewHandlerRegistry(
	authHandler *authhandler.AuthHandler,
	sessionHandler *sessionhandler.SessionHandler,
	userHandler *userhandler.UserHandler,
	projectHandler *projecthandler.ProjectHandler,
	appHandler *apphandler.AppHandler,
) *HandlerRegistry {
	return &HandlerRegistry{
		authHandler:    authHandler,
		sessionHandler: sessionHandler,
		userHandler:    userHandler,
		projectHandler: projectHandler,
		appHandler:     appHandler,
	}
}

func (s *HTTPServer) registerRoutes() {
	s.engine.GET("/", routeHome)
	s.engine.GET("/ping", routePing)
	s.engine.NoRoute(routeNotFound)

	// Swagger server
	if !s.config.IsProdEnv() {
		s.engine.StaticFile("/docs/openapi/swagger.json", "docs/openapi/swagger.json")
		s.engine.GET("/swagger/*any", swaggoGin.WrapHandler(swaggoFiles.Handler,
			swaggoGin.URL("/docs/openapi/swagger.json")))
	}

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

	{ // session group
		sessionGroup := apiV1Group.Group("/sessions")
		// User info
		sessionGroup.GET("/me", s.handlerRegistry.sessionHandler.GetMe)
		// Session handling
		sessionGroup.POST("/refresh", s.handlerRegistry.sessionHandler.RefreshSession)
		sessionGroup.DELETE("", s.handlerRegistry.sessionHandler.DeleteSession)
	}

	{ // auth group
		authGroup := apiV1Group.Group("/auth")
		// Login with password
		authGroup.POST("/login-with-password", s.handlerRegistry.sessionHandler.LoginWithPassword)
		authGroup.POST("/login-with-passcode", s.handlerRegistry.sessionHandler.LoginWithPasscode)
	}

	userGroup := apiV1Group.Group("/users")
	{ // user group
		// Get user info
		userGroup.GET("/base-list", s.handlerRegistry.userHandler.ListUserBase)
		userGroup.GET("/:userID", s.handlerRegistry.userHandler.GetUser)
		userGroup.GET("", s.handlerRegistry.userHandler.ListUser)
		// Password
		userGroup.PATCH("/current/password", s.handlerRegistry.userHandler.UpdateUserPassword)
	}

	projectGroup := apiV1Group.Group("/projects")
	{ // project group
		// Project info
		projectGroup.GET("/base-list", s.handlerRegistry.projectHandler.ListProjectBase)
		projectGroup.GET("/:projectID", s.handlerRegistry.projectHandler.GetProject)
		projectGroup.GET("", s.handlerRegistry.projectHandler.ListProject)
		// Creation & Update
		projectGroup.POST("", s.handlerRegistry.projectHandler.CreateProject)
		// Tags
		projectGroup.POST("/:projectID/tags", s.handlerRegistry.projectHandler.CreateProjectTag)
		projectGroup.POST("/:projectID/tags/delete", s.handlerRegistry.projectHandler.DeleteProjectTags)
		// Envs
		projectGroup.POST("/:projectID/envs", s.handlerRegistry.projectHandler.CreateProjectEnv)
		projectGroup.DELETE("/:projectID/envs/:projectEnvID", s.handlerRegistry.projectHandler.DeleteProjectEnv)
	}

	appGroup := apiV1Group.Group("/apps")
	{ // app group
		// App info
		appGroup.GET("/base-list", s.handlerRegistry.appHandler.ListAppBase)
		appGroup.GET("/:appID", s.handlerRegistry.appHandler.GetApp)
		appGroup.GET("", s.handlerRegistry.appHandler.ListApp)
	}
}

func routePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func routeNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, "not found")
}

func routeHome(c *gin.Context) {
	c.JSON(http.StatusOK, "localpaas api")
}
