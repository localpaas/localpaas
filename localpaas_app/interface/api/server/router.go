package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggoFiles "github.com/swaggo/files"
	swaggoGin "github.com/swaggo/gin-swagger"

	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/sessionhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/userhandler"
)

type HandlerRegistry struct {
	authHandler    *authhandler.AuthHandler
	sessionHandler *sessionhandler.SessionHandler
	userHandler    *userhandler.UserHandler
}

func NewHandlerRegistry(
	authHandler *authhandler.AuthHandler,
	sessionHandler *sessionhandler.SessionHandler,
	userHandler *userhandler.UserHandler,
) *HandlerRegistry {
	return &HandlerRegistry{
		authHandler:    authHandler,
		sessionHandler: sessionHandler,
		userHandler:    userHandler,
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
		userGroup.GET("/base-list", s.handlerRegistry.userHandler.ListUserSimple)
		userGroup.GET("/:userID", s.handlerRegistry.userHandler.GetUser)
		userGroup.GET("", s.handlerRegistry.userHandler.ListUser)
		// Password
		userGroup.PATCH("/current/password", s.handlerRegistry.userHandler.UpdateUserPassword)
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
