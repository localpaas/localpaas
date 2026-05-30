package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerWebhookRoutes(apiGroup *gin.RouterGroup) {
	webhookGroup := apiGroup.Group("/webhooks")
	webhookHandler := s.handlerRegistry.webhookHandler

	// App deployment
	webhookGroup.POST("/apps/:appToken/deploy", webhookHandler.WebhookDeployApp)

	// Repo webhook
	webhookGroup.POST("/:webhookID", webhookHandler.HandleRepoWebhook)
}
