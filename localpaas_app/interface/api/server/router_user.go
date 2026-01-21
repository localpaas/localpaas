package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerUserRoutes(apiGroup *gin.RouterGroup) (*gin.RouterGroup, *gin.RouterGroup) {
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

	return userGroup, userSettingGroup
}
