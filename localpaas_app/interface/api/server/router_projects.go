package server

import (
	"github.com/gin-gonic/gin"
)

//nolint:funlen
func (s *HTTPServer) registerProjectRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	projectGroup := apiGroup.Group("/projects")

	// Projects
	projectGroup.GET("/base", s.handlerRegistry.projectHandler.ListProjectBase)
	projectGroup.GET("/:projectID", s.handlerRegistry.projectHandler.GetProject)
	projectGroup.GET("", s.handlerRegistry.projectHandler.ListProject)
	projectGroup.POST("", s.handlerRegistry.projectHandler.CreateProject)
	projectGroup.DELETE("/:projectID", s.handlerRegistry.projectHandler.DeleteProject)

	// Tags
	projectGroup.POST("/:projectID/tags", s.handlerRegistry.projectHandler.CreateProjectTag)
	projectGroup.POST("/:projectID/tags/delete", s.handlerRegistry.projectHandler.DeleteProjectTags)

	// Env vars
	projectGroup.GET("/:projectID/env-vars", s.handlerRegistry.projectHandler.GetProjectEnvVars)
	projectGroup.PUT("/:projectID/env-vars", s.handlerRegistry.projectHandler.UpdateProjectEnvVars)

	// Secrets
	projectGroup.GET("/:projectID/secrets", s.handlerRegistry.projectHandler.ListProjectSecrets)
	projectGroup.POST("/:projectID/secrets", s.handlerRegistry.projectHandler.CreateProjectSecret)
	projectGroup.DELETE("/:projectID/secrets/:id", s.handlerRegistry.projectHandler.DeleteProjectSecret)

	// Cron jobs
	projectGroup.GET("/:projectID/cron-jobs", s.handlerRegistry.projectHandler.ListCronJob)
	projectGroup.GET("/:projectID/cron-jobs/:id", s.handlerRegistry.projectHandler.GetCronJob)
	projectGroup.POST("/:projectID/cron-jobs", s.handlerRegistry.projectHandler.CreateCronJob)
	projectGroup.PUT("/:projectID/cron-jobs/:id", s.handlerRegistry.projectHandler.UpdateCronJob)
	projectGroup.PUT("/:projectID/cron-jobs/:id/meta", s.handlerRegistry.projectHandler.UpdateCronJobMeta)
	projectGroup.DELETE("/:projectID/cron-jobs/:id", s.handlerRegistry.projectHandler.DeleteCronJob)

	// Settings import
	projectGroup.POST("/:projectID/settings-import", s.handlerRegistry.projectHandler.ImportSettings)

	// PROVIDERS
	projectProviderGroup := projectGroup.Group("/:projectID/providers")

	{ // github-app group
		githubAppGroup := projectProviderGroup.Group("/github-apps")
		// Info
		githubAppGroup.GET("/:id", s.handlerRegistry.projectHandler.GetGithubApp)
		githubAppGroup.GET("", s.handlerRegistry.projectHandler.ListGithubApp)
		// Creation & Update
		githubAppGroup.POST("", s.handlerRegistry.projectHandler.CreateGithubApp)
		githubAppGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateGithubApp)
		githubAppGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateGithubAppMeta)
		githubAppGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteGithubApp)
	}

	{ // git-token group
		gitTokenGroup := projectProviderGroup.Group("/git-tokens")
		// Info
		gitTokenGroup.GET("/:id", s.handlerRegistry.projectHandler.GetGitToken)
		gitTokenGroup.GET("", s.handlerRegistry.projectHandler.ListGitToken)
		// Creation & Update
		gitTokenGroup.POST("", s.handlerRegistry.projectHandler.CreateGitToken)
		gitTokenGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateGitToken)
		gitTokenGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateGitTokenMeta)
		gitTokenGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteGitToken)
	}

	{ // s3 storage group
		s3StorageGroup := projectProviderGroup.Group("/s3-storages")
		// Info
		s3StorageGroup.GET("/:id", s.handlerRegistry.projectHandler.GetS3Storage)
		s3StorageGroup.GET("", s.handlerRegistry.projectHandler.ListS3Storage)
		// Creation & Update
		s3StorageGroup.POST("", s.handlerRegistry.projectHandler.CreateS3Storage)
		s3StorageGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateS3Storage)
		s3StorageGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateS3StorageMeta)
		s3StorageGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteS3Storage)
	}

	{ // ssh key group
		sshKeyGroup := projectProviderGroup.Group("/ssh-keys")
		// Info
		sshKeyGroup.GET("/:id", s.handlerRegistry.projectHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.projectHandler.ListSSHKey)
		// Creation & Update
		sshKeyGroup.POST("", s.handlerRegistry.projectHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateSSHKey)
		sshKeyGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateSSHKeyMeta)
		sshKeyGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteSSHKey)
	}

	{ // slack group
		slackGroup := projectProviderGroup.Group("/slack")
		// Info
		slackGroup.GET("/:id", s.handlerRegistry.projectHandler.GetSlack)
		slackGroup.GET("", s.handlerRegistry.projectHandler.ListSlack)
		// Creation & Update
		slackGroup.POST("", s.handlerRegistry.projectHandler.CreateSlack)
		slackGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateSlack)
		slackGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateSlackMeta)
		slackGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteSlack)
	}

	{ // discord group
		discordGroup := projectProviderGroup.Group("/discord")
		// Info
		discordGroup.GET("/:id", s.handlerRegistry.projectHandler.GetDiscord)
		discordGroup.GET("", s.handlerRegistry.projectHandler.ListDiscord)
		// Creation & Update
		discordGroup.POST("", s.handlerRegistry.projectHandler.CreateDiscord)
		discordGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateDiscord)
		discordGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateDiscordMeta)
		discordGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteDiscord)
	}

	{ // registry auth group
		registryAuthGroup := projectProviderGroup.Group("/registry-auth")
		// Info
		registryAuthGroup.GET("/:id", s.handlerRegistry.projectHandler.GetRegistryAuth)
		registryAuthGroup.GET("", s.handlerRegistry.projectHandler.ListRegistryAuth)
		// Creation & Update
		registryAuthGroup.POST("", s.handlerRegistry.projectHandler.CreateRegistryAuth)
		registryAuthGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateRegistryAuth)
		registryAuthGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateRegistryAuthMeta)
		registryAuthGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteRegistryAuth)
	}

	{ // basic auth group
		basicAuthGroup := projectProviderGroup.Group("/basic-auth")
		// Info
		basicAuthGroup.GET("/:id", s.handlerRegistry.projectHandler.GetBasicAuth)
		basicAuthGroup.GET("", s.handlerRegistry.projectHandler.ListBasicAuth)
		// Creation & Update
		basicAuthGroup.POST("", s.handlerRegistry.projectHandler.CreateBasicAuth)
		basicAuthGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateBasicAuth)
		basicAuthGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateBasicAuthMeta)
		basicAuthGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteBasicAuth)
	}

	{ // ssl group
		sslGroup := projectProviderGroup.Group("/ssls")
		// Info
		sslGroup.GET("/:id", s.handlerRegistry.projectHandler.GetSsl)
		sslGroup.GET("", s.handlerRegistry.projectHandler.ListSsl)
		// Creation & Update
		sslGroup.POST("", s.handlerRegistry.projectHandler.CreateSsl)
		sslGroup.PUT("/:id", s.handlerRegistry.projectHandler.UpdateSsl)
		sslGroup.PUT("/:id/meta", s.handlerRegistry.projectHandler.UpdateSslMeta)
		sslGroup.DELETE("/:id", s.handlerRegistry.projectHandler.DeleteSsl)
	}

	return projectGroup
}
