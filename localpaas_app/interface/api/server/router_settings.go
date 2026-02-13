package server

import (
	"github.com/gin-gonic/gin"
)

//nolint:funlen
func (s *HTTPServer) registerSettingRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	settingGroup := apiGroup.Group("/settings")

	{ // oauth group
		oauthGroup := settingGroup.Group("/oauth")
		oauthGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetOAuth)
		oauthGroup.GET("", s.handlerRegistry.settingHandler.ListOAuth)
		oauthGroup.POST("", s.handlerRegistry.settingHandler.CreateOAuth)
		oauthGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateOAuth)
		oauthGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateOAuthMeta)
		oauthGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteOAuth)
	}

	{ // github-app group
		githubAppGroup := settingGroup.Group("/github-apps")
		githubAppGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetGithubApp)
		githubAppGroup.GET("", s.handlerRegistry.settingHandler.ListGithubApp)
		githubAppGroup.POST("", s.handlerRegistry.settingHandler.CreateGithubApp)
		githubAppGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateGithubApp)
		githubAppGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateGithubAppMeta)
		githubAppGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteGithubApp)
		// Test connection
		githubAppGroup.POST("/test-conn", s.handlerRegistry.settingHandler.TestGithubAppConn)
		// Installation
		githubAppGroup.POST("/installations/list", s.handlerRegistry.settingHandler.ListAppInstallation)
	}

	{ // access-token group
		accessTokenGroup := settingGroup.Group("/access-tokens")
		accessTokenGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetAccessToken)
		accessTokenGroup.GET("", s.handlerRegistry.settingHandler.ListAccessToken)
		accessTokenGroup.POST("", s.handlerRegistry.settingHandler.CreateAccessToken)
		accessTokenGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateAccessToken)
		accessTokenGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateAccessTokenMeta)
		accessTokenGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteAccessToken)
		// Test connection
		accessTokenGroup.POST("/test-conn", s.handlerRegistry.settingHandler.TestAccessTokenConn)
	}

	{ // aws group
		awsGroup := settingGroup.Group("/aws")
		awsGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetAWS)
		awsGroup.GET("", s.handlerRegistry.settingHandler.ListAWS)
		awsGroup.POST("", s.handlerRegistry.settingHandler.CreateAWS)
		awsGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateAWS)
		awsGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateAWSMeta)
		awsGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteAWS)
	}

	{ // aws s3 group
		awsS3Group := settingGroup.Group("/aws-s3")
		awsS3Group.GET("/:itemID", s.handlerRegistry.settingHandler.GetAWSS3)
		awsS3Group.GET("", s.handlerRegistry.settingHandler.ListAWSS3)
		awsS3Group.POST("", s.handlerRegistry.settingHandler.CreateAWSS3)
		awsS3Group.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateAWSS3)
		awsS3Group.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateAWSS3Meta)
		awsS3Group.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteAWSS3)
		// Test connection
		awsS3Group.POST("/test-conn", s.handlerRegistry.settingHandler.TestAWSS3Conn)
	}

	{ // ssh key group
		sshKeyGroup := settingGroup.Group("/ssh-keys")
		sshKeyGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.settingHandler.ListSSHKey)
		sshKeyGroup.POST("", s.handlerRegistry.settingHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateSSHKey)
		sshKeyGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateSSHKeyMeta)
		sshKeyGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteSSHKey)
	}

	{ // IM service group
		imServiceGroup := settingGroup.Group("/im-services")
		imServiceGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetIMService)
		imServiceGroup.GET("", s.handlerRegistry.settingHandler.ListIMService)
		imServiceGroup.POST("", s.handlerRegistry.settingHandler.CreateIMService)
		imServiceGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateIMService)
		imServiceGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateIMServiceMeta)
		imServiceGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteIMService)
		// Test connection
		imServiceGroup.POST("/test-send-msg", s.handlerRegistry.settingHandler.TestSendInstantMsg)
	}

	{ // registry auth group
		registryAuthGroup := settingGroup.Group("/registry-auth")
		registryAuthGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetRegistryAuth)
		registryAuthGroup.GET("", s.handlerRegistry.settingHandler.ListRegistryAuth)
		registryAuthGroup.POST("", s.handlerRegistry.settingHandler.CreateRegistryAuth)
		registryAuthGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateRegistryAuth)
		registryAuthGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateRegistryAuthMeta)
		registryAuthGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteRegistryAuth)
		// Test connection
		registryAuthGroup.POST("/test-conn", s.handlerRegistry.settingHandler.TestRegistryAuthConn)
	}

	{ // basic auth group
		basicAuthGroup := settingGroup.Group("/basic-auth")
		basicAuthGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetBasicAuth)
		basicAuthGroup.GET("", s.handlerRegistry.settingHandler.ListBasicAuth)
		basicAuthGroup.POST("", s.handlerRegistry.settingHandler.CreateBasicAuth)
		basicAuthGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateBasicAuth)
		basicAuthGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateBasicAuthMeta)
		basicAuthGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteBasicAuth)
	}

	{ // ssl group
		sslGroup := settingGroup.Group("/ssls")
		sslGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetSSL)
		sslGroup.GET("", s.handlerRegistry.settingHandler.ListSSL)
		sslGroup.POST("", s.handlerRegistry.settingHandler.CreateSSL)
		sslGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateSSL)
		sslGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateSSLMeta)
		sslGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteSSL)
	}

	{ // email group
		emailGroup := settingGroup.Group("/emails")
		emailGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetEmail)
		emailGroup.GET("", s.handlerRegistry.settingHandler.ListEmail)
		emailGroup.POST("", s.handlerRegistry.settingHandler.CreateEmail)
		emailGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateEmail)
		emailGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateEmailMeta)
		emailGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteEmail)
		// Test connection
		emailGroup.POST("/test-send-mail", s.handlerRegistry.settingHandler.TestSendMail)
	}

	{ // secrets group
		secretGroup := settingGroup.Group("/secrets")
		secretGroup.GET("", s.handlerRegistry.settingHandler.ListSecret)
		secretGroup.POST("", s.handlerRegistry.settingHandler.CreateSecret)
		secretGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateSecret)
		secretGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateSecretMeta)
		secretGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteSecret)
	}

	{ // cron-job group
		cronJobGroup := settingGroup.Group("/cron-jobs")
		cronJobGroup.GET("/:itemID", s.handlerRegistry.settingHandler.GetCronJob)
		cronJobGroup.GET("", s.handlerRegistry.settingHandler.ListCronJob)
		cronJobGroup.POST("", s.handlerRegistry.settingHandler.CreateCronJob)
		cronJobGroup.PUT("/:itemID", s.handlerRegistry.settingHandler.UpdateCronJob)
		cronJobGroup.PUT("/:itemID/meta", s.handlerRegistry.settingHandler.UpdateCronJobMeta)
		cronJobGroup.DELETE("/:itemID", s.handlerRegistry.settingHandler.DeleteCronJob)
	}

	return settingGroup
}
