package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerSettingRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	settingGroup := apiGroup.Group("/settings")

	{ // secrets group
		secretGroup := settingGroup.Group("/secrets")
		// Info
		secretGroup.GET("", s.handlerRegistry.settingHandler.ListSecret)
		// Creation & Update
		secretGroup.POST("", s.handlerRegistry.settingHandler.CreateSecret)
		secretGroup.PUT("/:id/meta", s.handlerRegistry.settingHandler.UpdateSecretMeta)
		secretGroup.DELETE("/:id", s.handlerRegistry.settingHandler.DeleteSecret)
	}

	{ // cron-job group
		cronJobGroup := settingGroup.Group("/cron-jobs")
		// Info
		cronJobGroup.GET("/:id", s.handlerRegistry.settingHandler.GetCronJob)
		cronJobGroup.GET("", s.handlerRegistry.settingHandler.ListCronJob)
		// Creation & Update
		cronJobGroup.POST("", s.handlerRegistry.settingHandler.CreateCronJob)
		cronJobGroup.PUT("/:id", s.handlerRegistry.settingHandler.UpdateCronJob)
		cronJobGroup.PUT("/:id/meta", s.handlerRegistry.settingHandler.UpdateCronJobMeta)
		cronJobGroup.DELETE("/:id", s.handlerRegistry.settingHandler.DeleteCronJob)
	}

	return settingGroup
}
