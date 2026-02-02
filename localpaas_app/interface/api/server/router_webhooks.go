package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerWebhookRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	webhookGroup := apiGroup.Group("/webhooks")

	webhookGroup.POST("/apps/:appToken/deploy", s.handlerRegistry.webhookHandler.WebhookDeployApp)

	// Github
	webhookGroup.POST("/github", s.handlerRegistry.webhookHandler.HandleGithubWebhook)
	// Gitlab
	webhookGroup.POST("/gitlab", s.handlerRegistry.webhookHandler.HandleGitlabWebhook)
	// Gitea
	webhookGroup.POST("/gitea", s.handlerRegistry.webhookHandler.HandleGiteaWebhook)
	// Bitbucket
	webhookGroup.POST("/bitbucket", s.handlerRegistry.webhookHandler.HandleBitbucketWebhook)

	return webhookGroup
}
