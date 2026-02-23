package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerWebhookRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	webhookGroup := apiGroup.Group("/webhooks")
	webhookHandler := s.handlerRegistry.webhookHandler

	// App deployment
	webhookGroup.POST("/apps/:appToken/deploy", webhookHandler.WebhookDeployApp)

	// Repo webhook
	webhookGroup.POST("/:kind/:secret", webhookHandler.HandleRepoWebhook)

	return webhookGroup
}
