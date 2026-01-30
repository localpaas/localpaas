package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerAppRoutes(projectGroup *gin.RouterGroup) *gin.RouterGroup {
	appGroup := projectGroup.Group("/:projectID/apps")

	// Info
	appGroup.GET("/base", s.handlerRegistry.appHandler.ListAppBase)
	appGroup.GET("/:appID", s.handlerRegistry.appHandler.GetApp)
	appGroup.GET("", s.handlerRegistry.appHandler.ListApp)
	// Creation & Update
	appGroup.POST("", s.handlerRegistry.appHandler.CreateApp)
	appGroup.PUT("/:appID", s.handlerRegistry.appHandler.UpdateApp)
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
	appGroup.GET("/:appID/secrets", s.handlerRegistry.appHandler.ListAppSecret)
	appGroup.POST("/:appID/secrets", s.handlerRegistry.appHandler.CreateAppSecret)
	appGroup.PUT("/:appID/secrets/:id", s.handlerRegistry.appHandler.UpdateAppSecret)
	appGroup.DELETE("/:appID/secrets/:id", s.handlerRegistry.appHandler.DeleteAppSecret)
	// Domain SSL
	appGroup.POST("/:appID/ssl/obtain", s.handlerRegistry.appHandler.ObtainDomainSSL)
	// Logs
	appGroup.GET("/:appID/runtime-logs", func(ctx *gin.Context) {
		s.handlerRegistry.appHandler.GetAppRuntimeLogs(ctx, s.websocket)
	})

	appDeploymentGroup := appGroup.Group("/:appID/deployments")
	{ // app deployment group
		// Info
		appDeploymentGroup.GET("/:id", s.handlerRegistry.appHandler.GetAppDeployment)
		appDeploymentGroup.GET("", s.handlerRegistry.appHandler.ListAppDeployment)
		// Cancel
		appDeploymentGroup.POST("/:id/cancel", s.handlerRegistry.appHandler.CancelAppDeployment)
		// Logs
		appDeploymentGroup.GET("/:id/logs", func(ctx *gin.Context) {
			s.handlerRegistry.appHandler.GetAppDeploymentLogs(ctx, s.websocket)
		})
	}

	return appGroup
}
