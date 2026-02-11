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
	// Token
	appGroup.PUT("/:appID/token", s.handlerRegistry.appHandler.UpdateAppToken)
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
	appGroup.PUT("/:appID/secrets/:itemID", s.handlerRegistry.appHandler.UpdateAppSecret)
	appGroup.DELETE("/:appID/secrets/:itemID", s.handlerRegistry.appHandler.DeleteAppSecret)
	// Domain SSL
	appGroup.POST("/:appID/ssl/obtain", s.handlerRegistry.appHandler.ObtainDomainSSL)
	// Logs
	appGroup.GET("/:appID/runtime-logs", func(ctx *gin.Context) {
		s.handlerRegistry.appHandler.GetAppRuntimeLogs(ctx, s.websocket)
	})

	cronJobGroup := appGroup.Group("/:appID/cron-jobs")
	{ // Cron job group
		cronJobGroup.GET("", s.handlerRegistry.appHandler.ListAppCronJob)
		cronJobGroup.GET("/:itemID", s.handlerRegistry.appHandler.GetAppCronJob)
		cronJobGroup.POST("", s.handlerRegistry.appHandler.CreateAppCronJob)
		cronJobGroup.PUT("/:itemID", s.handlerRegistry.appHandler.UpdateAppCronJob)
		cronJobGroup.PUT("/:itemID/meta", s.handlerRegistry.appHandler.UpdateAppCronJobMeta)
		cronJobGroup.DELETE("/:itemID", s.handlerRegistry.appHandler.DeleteAppCronJob)
		// Execute
		cronJobGroup.POST("/:itemID/exec", s.handlerRegistry.appHandler.ExecuteAppCronJob)

		// Cron job task group
		cronJobGroup.GET("/:itemID/tasks", s.handlerRegistry.appHandler.ListAppCronJobTask)
		cronJobGroup.GET("/:itemID/tasks/:taskID", s.handlerRegistry.appHandler.GetAppCronJobTask)
		cronJobGroup.GET("/:itemID/tasks/:taskID/logs", func(ctx *gin.Context) {
			s.handlerRegistry.appHandler.GetAppCronJobTaskLogs(ctx, s.websocket)
		})
	}

	healthcheckGroup := appGroup.Group("/:appID/healthchecks")
	{ // healthcheck group
		healthcheckGroup.GET("", s.handlerRegistry.appHandler.ListAppHealthcheck)
		healthcheckGroup.GET("/:itemID", s.handlerRegistry.appHandler.GetAppHealthcheck)
		healthcheckGroup.POST("", s.handlerRegistry.appHandler.CreateAppHealthcheck)
		healthcheckGroup.PUT("/:itemID", s.handlerRegistry.appHandler.UpdateAppHealthcheck)
		healthcheckGroup.PUT("/:itemID/meta", s.handlerRegistry.appHandler.UpdateAppHealthcheckMeta)
		healthcheckGroup.DELETE("/:itemID", s.handlerRegistry.appHandler.DeleteAppHealthcheck)

		// Healthcheck task group
		healthcheckGroup.GET("/:itemID/tasks", s.handlerRegistry.appHandler.ListAppHealthcheckTask)
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

	return appGroup
}
