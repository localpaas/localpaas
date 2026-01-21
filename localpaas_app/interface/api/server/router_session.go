package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerSessionRoutes(apiGroup *gin.RouterGroup) (*gin.RouterGroup, *gin.RouterGroup) {
	sessionGroup := apiGroup.Group("/sessions")
	{
		// User info
		sessionGroup.GET("/me", s.handlerRegistry.sessionHandler.GetMe)
		// Session handling
		sessionGroup.POST("/refresh", s.handlerRegistry.sessionHandler.RefreshSession)
		sessionGroup.DELETE("", s.handlerRegistry.sessionHandler.DeleteSession)
		sessionGroup.POST("/delete-all", s.handlerRegistry.sessionHandler.DeleteAllSessions)
	}

	authGroup := apiGroup.Group("/auth")
	{
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

	return sessionGroup, authGroup
}
