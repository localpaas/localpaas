package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerSystemRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	systemGroup := apiGroup.Group("/system")

	{ // task group
		taskGroup := systemGroup.Group("/tasks")
		taskGroup.GET("", s.handlerRegistry.systemHandler.ListTask)
		taskGroup.GET("/:taskID", s.handlerRegistry.systemHandler.GetTask)
		taskGroup.PUT("/:taskID/meta", s.handlerRegistry.systemHandler.UpdateTaskMeta)
		taskGroup.POST("/:taskID/cancel", s.handlerRegistry.systemHandler.CancelTask)
	}

	{ // error group
		errorGroup := systemGroup.Group("/errors")
		errorGroup.GET("", s.handlerRegistry.systemHandler.ListSysError)
		errorGroup.GET("/:errorID", s.handlerRegistry.systemHandler.GetSysError)
		errorGroup.DELETE("/:errorID", s.handlerRegistry.systemHandler.DeleteSysError)
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
