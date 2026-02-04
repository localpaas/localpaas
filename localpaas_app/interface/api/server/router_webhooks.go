package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerWebhookRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	webhookGroup := apiGroup.Group("/webhooks")

	// App deployment
	webhookGroup.POST("/apps/:appToken/deploy", s.handlerRegistry.webhookHandler.WebhookDeployApp)

	// Repo webhook
	webhookGroup.POST("/:kind/:secret", s.handlerRegistry.webhookHandler.HandleRepoWebhook)

	return webhookGroup
}
