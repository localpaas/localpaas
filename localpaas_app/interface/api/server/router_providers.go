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

	{ // aws group
		awsGroup := providerGroup.Group("/aws")
		// Info
		awsGroup.GET("/:id", s.handlerRegistry.providersHandler.GetAWS)
		awsGroup.GET("", s.handlerRegistry.providersHandler.ListAWS)
		// Creation & Update
		awsGroup.POST("", s.handlerRegistry.providersHandler.CreateAWS)
		awsGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateAWS)
		awsGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateAWSMeta)
		awsGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteAWS)
	}

	{ // aws s3 group
		awsS3Group := providerGroup.Group("/aws-s3")
		// Info
		awsS3Group.GET("/:id", s.handlerRegistry.providersHandler.GetAWSS3)
		awsS3Group.GET("", s.handlerRegistry.providersHandler.ListAWSS3)
		// Creation & Update
		awsS3Group.POST("", s.handlerRegistry.providersHandler.CreateAWSS3)
		awsS3Group.PUT("/:id", s.handlerRegistry.providersHandler.UpdateAWSS3)
		awsS3Group.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateAWSS3Meta)
		awsS3Group.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteAWSS3)
		// Test connection
		awsS3Group.POST("/test-conn", s.handlerRegistry.providersHandler.TestAWSS3Conn)
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

	{ // IM service group
		imServiceGroup := providerGroup.Group("/im-services")
		// Info
		imServiceGroup.GET("/:id", s.handlerRegistry.providersHandler.GetIMService)
		imServiceGroup.GET("", s.handlerRegistry.providersHandler.ListIMService)
		// Creation & Update
		imServiceGroup.POST("", s.handlerRegistry.providersHandler.CreateIMService)
		imServiceGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateIMService)
		imServiceGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateIMServiceMeta)
		imServiceGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteIMService)
		// Test connection
		imServiceGroup.POST("/test-send-msg", s.handlerRegistry.providersHandler.TestSendInstantMsg)
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

	{ // email group
		emailGroup := providerGroup.Group("/emails")
		// Info
		emailGroup.GET("/:id", s.handlerRegistry.providersHandler.GetEmail)
		emailGroup.GET("", s.handlerRegistry.providersHandler.ListEmail)
		// Creation & Update
		emailGroup.POST("", s.handlerRegistry.providersHandler.CreateEmail)
		emailGroup.PUT("/:id", s.handlerRegistry.providersHandler.UpdateEmail)
		emailGroup.PUT("/:id/meta", s.handlerRegistry.providersHandler.UpdateEmailMeta)
		emailGroup.DELETE("/:id", s.handlerRegistry.providersHandler.DeleteEmail)
		// Test connection
		emailGroup.POST("/test-send-mail", s.handlerRegistry.providersHandler.TestSendMail)
	}

	return providerGroup
}
