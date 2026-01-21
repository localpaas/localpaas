package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerProviderRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	providerGroup := apiGroup.Group("/providers")

	{ // oauth group
		oauthGroup := providerGroup.Group("/oauth")
		// Info
		oauthGroup.GET("/:id", s.handlerRegistry.providersHandler.GetOAuth)
		oauthGroup.GET("", s.handlerRegistry.providersHandler.ListOAuth)
		// Creation & Update
		oauthGroup.POST("", s.handlerRegistry.providersHandler.CreateOAuth)
		oauthGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateOAuth)
		oauthGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateOAuthMeta)
		oauthGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteOAuth)
	}

	{ // github-app group
		githubAppGroup := providerGroup.Group("/github-apps")
		// Info
		githubAppGroup.GET("/:id", s.handlerRegistry.providersHandler.GetGithubApp)
		githubAppGroup.GET("", s.handlerRegistry.providersHandler.ListGithubApp)
		// Creation & Update
		githubAppGroup.POST("", s.handlerRegistry.providersHandler.CreateGithubApp)
		githubAppGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateGithubApp)
		githubAppGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateGithubAppMeta)
		githubAppGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteGithubApp)
		// Test connection
		githubAppGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestGithubAppConn)
		// Installation
		githubAppGroup.POST("/installations/list", s.handlerRegistry.providersHandler.ListAppInstallation)
	}

	{ // git-token group
		gitTokenGroup := providerGroup.Group("/git-tokens")
		// Info
		gitTokenGroup.GET("/:id", s.handlerRegistry.providersHandler.GetGitToken)
		gitTokenGroup.GET("", s.handlerRegistry.providersHandler.ListGitToken)
		// Creation & Update
		gitTokenGroup.POST("", s.handlerRegistry.providersHandler.CreateGitToken)
		gitTokenGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateGitToken)
		gitTokenGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateGitTokenMeta)
		gitTokenGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteGitToken)
		// Test connection
		gitTokenGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestGitTokenConn)
	}

	{ // s3 storage group
		s3StorageGroup := providerGroup.Group("/s3-storages")
		// Info
		s3StorageGroup.GET("/:id", s.handlerRegistry.providersHandler.GetS3Storage)
		s3StorageGroup.GET("", s.handlerRegistry.providersHandler.ListS3Storage)
		// Creation & Update
		s3StorageGroup.POST("", s.handlerRegistry.providersHandler.CreateS3Storage)
		s3StorageGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateS3Storage)
		s3StorageGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateS3StorageMeta)
		s3StorageGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteS3Storage)
		// Test connection
		s3StorageGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestS3StorageConn)
	}

	{ // ssh key group
		sshKeyGroup := providerGroup.Group("/ssh-keys")
		// Info
		sshKeyGroup.GET("/:id", s.handlerRegistry.providersHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.providersHandler.ListSSHKey)
		// Creation & Update
		sshKeyGroup.POST("", s.handlerRegistry.providersHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateSSHKey)
		sshKeyGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateSSHKeyMeta)
		sshKeyGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteSSHKey)
	}

	{ // slack group
		slackGroup := providerGroup.Group("/slack")
		// Info
		slackGroup.GET("/:id", s.handlerRegistry.providersHandler.GetSlack)
		slackGroup.GET("", s.handlerRegistry.providersHandler.ListSlack)
		// Creation & Update
		slackGroup.POST("", s.handlerRegistry.providersHandler.CreateSlack)
		slackGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateSlack)
		slackGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateSlackMeta)
		slackGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteSlack)
		// Test connection
		slackGroup.POST("/test-send-msg", s.handlerRegistry.providersHandler.TestSendSlackMsg)
	}

	{ // discord group
		discordGroup := providerGroup.Group("/discord")
		// Info
		discordGroup.GET("/:id", s.handlerRegistry.providersHandler.GetDiscord)
		discordGroup.GET("", s.handlerRegistry.providersHandler.ListDiscord)
		// Creation & Update
		discordGroup.POST("", s.handlerRegistry.providersHandler.CreateDiscord)
		discordGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateDiscord)
		discordGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateDiscordMeta)
		discordGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteDiscord)
		// Test connection
		discordGroup.POST("/test-send-msg", s.handlerRegistry.providersHandler.TestSendDiscordMsg)
	}

	{ // registry auth group
		registryAuthGroup := providerGroup.Group("/registry-auth")
		// Info
		registryAuthGroup.GET("/:id", s.handlerRegistry.providersHandler.GetRegistryAuth)
		registryAuthGroup.GET("", s.handlerRegistry.providersHandler.ListRegistryAuth)
		// Creation & Update
		registryAuthGroup.POST("", s.handlerRegistry.providersHandler.CreateRegistryAuth)
		registryAuthGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateRegistryAuth)
		registryAuthGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateRegistryAuthMeta)
		registryAuthGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteRegistryAuth)
		// Test connection
		registryAuthGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestRegistryAuthConn)
	}

	{ // basic auth group
		basicAuthGroup := providerGroup.Group("/basic-auth")
		// Info
		basicAuthGroup.GET("/:id", s.handlerRegistry.providersHandler.GetBasicAuth)
		basicAuthGroup.GET("", s.handlerRegistry.providersHandler.ListBasicAuth)
		// Creation & Update
		basicAuthGroup.POST("", s.handlerRegistry.providersHandler.CreateBasicAuth)
		basicAuthGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateBasicAuth)
		basicAuthGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateBasicAuthMeta)
		basicAuthGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteBasicAuth)
	}

	{ // ssl group
		sslGroup := providerGroup.Group("/ssls")
		// Info
		sslGroup.GET("/:id", s.handlerRegistry.providersHandler.GetSsl)
		sslGroup.GET("", s.handlerRegistry.providersHandler.ListSsl)
		// Creation & Update
		sslGroup.POST("", s.handlerRegistry.providersHandler.CreateSsl)
		sslGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateSsl)
		sslGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateSslMeta)
		sslGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteSsl)
	}

	return providerGroup
}
