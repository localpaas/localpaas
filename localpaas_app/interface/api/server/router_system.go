package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerSystemRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
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

	return systemGroup
}
