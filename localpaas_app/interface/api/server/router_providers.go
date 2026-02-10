package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerProviderRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	providerGroup := apiGroup.Group("/providers")

	{ // oauth group
		oauthGroup := providerGroup.Group("/oauth")
		oauthGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetOAuth)
		oauthGroup.GET("", s.handlerRegistry.providersHandler.ListOAuth)
		oauthGroup.POST("", s.handlerRegistry.providersHandler.CreateOAuth)
		oauthGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateOAuth)
		oauthGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateOAuthMeta)
		oauthGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteOAuth)
	}

	{ // github-app group
		githubAppGroup := providerGroup.Group("/github-apps")
		githubAppGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetGithubApp)
		githubAppGroup.GET("", s.handlerRegistry.providersHandler.ListGithubApp)
		githubAppGroup.POST("", s.handlerRegistry.providersHandler.CreateGithubApp)
		githubAppGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateGithubApp)
		githubAppGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateGithubAppMeta)
		githubAppGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteGithubApp)
		// Test connection
		githubAppGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestGithubAppConn)
		// Installation
		githubAppGroup.POST("/installations/list", s.handlerRegistry.providersHandler.ListAppInstallation)
	}

	{ // access-token group
		accessTokenGroup := providerGroup.Group("/access-tokens")
		accessTokenGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetAccessToken)
		accessTokenGroup.GET("", s.handlerRegistry.providersHandler.ListAccessToken)
		accessTokenGroup.POST("", s.handlerRegistry.providersHandler.CreateAccessToken)
		accessTokenGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateAccessToken)
		accessTokenGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateAccessTokenMeta)
		accessTokenGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteAccessToken)
		// Test connection
		accessTokenGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestAccessTokenConn)
	}

	{ // aws group
		awsGroup := providerGroup.Group("/aws")
		awsGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetAWS)
		awsGroup.GET("", s.handlerRegistry.providersHandler.ListAWS)
		awsGroup.POST("", s.handlerRegistry.providersHandler.CreateAWS)
		awsGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateAWS)
		awsGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateAWSMeta)
		awsGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteAWS)
	}

	{ // aws s3 group
		awsS3Group := providerGroup.Group("/aws-s3")
		awsS3Group.GET("/:settingID", s.handlerRegistry.providersHandler.GetAWSS3)
		awsS3Group.GET("", s.handlerRegistry.providersHandler.ListAWSS3)
		awsS3Group.POST("", s.handlerRegistry.providersHandler.CreateAWSS3)
		awsS3Group.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateAWSS3)
		awsS3Group.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateAWSS3Meta)
		awsS3Group.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteAWSS3)
		// Test connection
		awsS3Group.POST("/test-conn", s.handlerRegistry.providersHandler.TestAWSS3Conn)
	}

	{ // ssh key group
		sshKeyGroup := providerGroup.Group("/ssh-keys")
		sshKeyGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.providersHandler.ListSSHKey)
		sshKeyGroup.POST("", s.handlerRegistry.providersHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateSSHKey)
		sshKeyGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateSSHKeyMeta)
		sshKeyGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteSSHKey)
	}

	{ // IM service group
		imServiceGroup := providerGroup.Group("/im-services")
		imServiceGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetIMService)
		imServiceGroup.GET("", s.handlerRegistry.providersHandler.ListIMService)
		imServiceGroup.POST("", s.handlerRegistry.providersHandler.CreateIMService)
		imServiceGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateIMService)
		imServiceGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateIMServiceMeta)
		imServiceGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteIMService)
		// Test connection
		imServiceGroup.POST("/test-send-msg", s.handlerRegistry.providersHandler.TestSendInstantMsg)
	}

	{ // registry auth group
		registryAuthGroup := providerGroup.Group("/registry-auth")
		registryAuthGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetRegistryAuth)
		registryAuthGroup.GET("", s.handlerRegistry.providersHandler.ListRegistryAuth)
		registryAuthGroup.POST("", s.handlerRegistry.providersHandler.CreateRegistryAuth)
		registryAuthGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateRegistryAuth)
		registryAuthGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateRegistryAuthMeta)
		registryAuthGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteRegistryAuth)
		// Test connection
		registryAuthGroup.POST("/test-conn", s.handlerRegistry.providersHandler.TestRegistryAuthConn)
	}

	{ // basic auth group
		basicAuthGroup := providerGroup.Group("/basic-auth")
		basicAuthGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetBasicAuth)
		basicAuthGroup.GET("", s.handlerRegistry.providersHandler.ListBasicAuth)
		basicAuthGroup.POST("", s.handlerRegistry.providersHandler.CreateBasicAuth)
		basicAuthGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateBasicAuth)
		basicAuthGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateBasicAuthMeta)
		basicAuthGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteBasicAuth)
	}

	{ // ssl group
		sslGroup := providerGroup.Group("/ssls")
		sslGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetSSL)
		sslGroup.GET("", s.handlerRegistry.providersHandler.ListSSL)
		sslGroup.POST("", s.handlerRegistry.providersHandler.CreateSSL)
		sslGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateSSL)
		sslGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateSSLMeta)
		sslGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteSSL)
	}

	{ // email group
		emailGroup := providerGroup.Group("/emails")
		emailGroup.GET("/:settingID", s.handlerRegistry.providersHandler.GetEmail)
		emailGroup.GET("", s.handlerRegistry.providersHandler.ListEmail)
		emailGroup.POST("", s.handlerRegistry.providersHandler.CreateEmail)
		emailGroup.PUT("/:settingID", s.handlerRegistry.providersHandler.UpdateEmail)
		emailGroup.PUT("/:settingID/meta", s.handlerRegistry.providersHandler.UpdateEmailMeta)
		emailGroup.DELETE("/:settingID", s.handlerRegistry.providersHandler.DeleteEmail)
		// Test connection
		emailGroup.POST("/test-send-mail", s.handlerRegistry.providersHandler.TestSendMail)
	}

	return providerGroup
}
