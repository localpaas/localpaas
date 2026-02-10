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
	projectGroup.PUT("/:projectID", s.handlerRegistry.projectHandler.UpdateProject)
	projectGroup.PUT("/:projectID/photo", s.handlerRegistry.projectHandler.UpdateProjectPhoto)
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
	projectGroup.PUT("/:projectID/secrets/:itemID", s.handlerRegistry.projectHandler.UpdateProjectSecret)
	projectGroup.DELETE("/:projectID/secrets/:itemID", s.handlerRegistry.projectHandler.DeleteProjectSecret)

	// Cron jobs
	projectGroup.GET("/:projectID/cron-jobs", s.handlerRegistry.projectHandler.ListCronJob)
	projectGroup.GET("/:projectID/cron-jobs/:itemID", s.handlerRegistry.projectHandler.GetCronJob)
	projectGroup.POST("/:projectID/cron-jobs", s.handlerRegistry.projectHandler.CreateCronJob)
	projectGroup.PUT("/:projectID/cron-jobs/:itemID", s.handlerRegistry.projectHandler.UpdateCronJob)
	projectGroup.PUT("/:projectID/cron-jobs/:itemID/meta", s.handlerRegistry.projectHandler.UpdateCronJobMeta)
	projectGroup.DELETE("/:projectID/cron-jobs/:itemID", s.handlerRegistry.projectHandler.DeleteCronJob)

	// Settings import
	projectGroup.POST("/:projectID/settings-import", s.handlerRegistry.projectHandler.ImportSettings)

	// PROVIDERS
	projectProviderGroup := projectGroup.Group("/:projectID/providers")

	{ // github-app group
		githubAppGroup := projectProviderGroup.Group("/github-apps")
		githubAppGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetGithubApp)
		githubAppGroup.GET("", s.handlerRegistry.projectHandler.ListGithubApp)
		githubAppGroup.POST("", s.handlerRegistry.projectHandler.CreateGithubApp)
		githubAppGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateGithubApp)
		githubAppGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateGithubAppMeta)
		githubAppGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteGithubApp)
	}

	{ // access-token group
		accessTokenGroup := projectProviderGroup.Group("/access-tokens")
		accessTokenGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetAccessToken)
		accessTokenGroup.GET("", s.handlerRegistry.projectHandler.ListAccessToken)
		accessTokenGroup.POST("", s.handlerRegistry.projectHandler.CreateAccessToken)
		accessTokenGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateAccessToken)
		accessTokenGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateAccessTokenMeta)
		accessTokenGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteAccessToken)
	}

	{ // aws group
		awsGroup := projectProviderGroup.Group("/aws")
		awsGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetAWS)
		awsGroup.GET("", s.handlerRegistry.projectHandler.ListAWS)
		awsGroup.POST("", s.handlerRegistry.projectHandler.CreateAWS)
		awsGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateAWS)
		awsGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateAWSMeta)
		awsGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteAWS)
	}

	{ // aws s3 group
		awsS3Group := projectProviderGroup.Group("/aws-s3")
		awsS3Group.GET("/:itemID", s.handlerRegistry.projectHandler.GetAWSS3)
		awsS3Group.GET("", s.handlerRegistry.projectHandler.ListAWSS3)
		awsS3Group.POST("", s.handlerRegistry.projectHandler.CreateAWSS3)
		awsS3Group.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateAWSS3)
		awsS3Group.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateAWSS3Meta)
		awsS3Group.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteAWSS3)
	}

	{ // ssh key group
		sshKeyGroup := projectProviderGroup.Group("/ssh-keys")
		sshKeyGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.projectHandler.ListSSHKey)
		sshKeyGroup.POST("", s.handlerRegistry.projectHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateSSHKey)
		sshKeyGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateSSHKeyMeta)
		sshKeyGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteSSHKey)
	}

	{ // IM service group
		imServiceGroup := projectProviderGroup.Group("/im-services")
		imServiceGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetIMService)
		imServiceGroup.GET("", s.handlerRegistry.projectHandler.ListIMService)
		imServiceGroup.POST("", s.handlerRegistry.projectHandler.CreateIMService)
		imServiceGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateIMService)
		imServiceGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateIMServiceMeta)
		imServiceGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteIMService)
	}

	{ // registry auth group
		registryAuthGroup := projectProviderGroup.Group("/registry-auth")
		registryAuthGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetRegistryAuth)
		registryAuthGroup.GET("", s.handlerRegistry.projectHandler.ListRegistryAuth)
		registryAuthGroup.POST("", s.handlerRegistry.projectHandler.CreateRegistryAuth)
		registryAuthGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateRegistryAuth)
		registryAuthGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateRegistryAuthMeta)
		registryAuthGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteRegistryAuth)
	}

	{ // basic auth group
		basicAuthGroup := projectProviderGroup.Group("/basic-auth")
		basicAuthGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetBasicAuth)
		basicAuthGroup.GET("", s.handlerRegistry.projectHandler.ListBasicAuth)
		basicAuthGroup.POST("", s.handlerRegistry.projectHandler.CreateBasicAuth)
		basicAuthGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateBasicAuth)
		basicAuthGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateBasicAuthMeta)
		basicAuthGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteBasicAuth)
	}

	{ // ssl group
		sslGroup := projectProviderGroup.Group("/ssls")
		sslGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetSSL)
		sslGroup.GET("", s.handlerRegistry.projectHandler.ListSSL)
		sslGroup.POST("", s.handlerRegistry.projectHandler.CreateSSL)
		sslGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateSSL)
		sslGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateSSLMeta)
		sslGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteSSL)
	}

	{ // email group
		emailGroup := projectProviderGroup.Group("/emails")
		emailGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetEmail)
		emailGroup.GET("", s.handlerRegistry.projectHandler.ListEmail)
		emailGroup.POST("", s.handlerRegistry.projectHandler.CreateEmail)
		emailGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateEmail)
		emailGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateEmailMeta)
		emailGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteEmail)
	}

	{ // repo webhook group
		repoWebhookGroup := projectProviderGroup.Group("/repo-webhooks")
		repoWebhookGroup.GET("/:itemID", s.handlerRegistry.projectHandler.GetRepoWebhook)
		repoWebhookGroup.GET("", s.handlerRegistry.projectHandler.ListRepoWebhook)
		repoWebhookGroup.POST("", s.handlerRegistry.projectHandler.CreateRepoWebhook)
		repoWebhookGroup.PUT("/:itemID", s.handlerRegistry.projectHandler.UpdateRepoWebhook)
		repoWebhookGroup.PUT("/:itemID/meta", s.handlerRegistry.projectHandler.UpdateRepoWebhookMeta)
		repoWebhookGroup.DELETE("/:itemID", s.handlerRegistry.projectHandler.DeleteRepoWebhook)
	}

	return projectGroup
}
