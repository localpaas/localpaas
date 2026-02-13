package server

import (
	"github.com/gin-gonic/gin"
)

//nolint:funlen
func (s *HTTPServer) registerSettingRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	settingGroup := apiGroup.Group("/settings")

	{ // oauth group
		oauthGroup := settingGroup.Group("/oauth")
		oauthGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetOAuth)
		oauthGroup.GET("", s.handlerRegistry.settingHandler.ListOAuth)
		oauthGroup.POST("", s.handlerRegistry.settingHandler.CreateOAuth)
		oauthGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateOAuth)
		oauthGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateOAuthMeta)
		oauthGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteOAuth)
	}

	{ // github-app group
		githubAppGroup := settingGroup.Group("/github-apps")
		githubAppGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetGithubApp)
		githubAppGroup.GET("", s.handlerRegistry.settingHandler.ListGithubApp)
		githubAppGroup.POST("", s.handlerRegistry.settingHandler.CreateGithubApp)
		githubAppGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateGithubApp)
		githubAppGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateGithubAppMeta)
		githubAppGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteGithubApp)
		// Test connection
		githubAppGroup.POST("/test-conn", s.handlerRegistry.settingHandler.TestGithubAppConn)
		// Installation
		githubAppGroup.POST("/installations/list", s.handlerRegistry.settingHandler.ListAppInstallation)
	}

	{ // access-token group
		accessTokenGroup := settingGroup.Group("/access-tokens")
		accessTokenGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetAccessToken)
		accessTokenGroup.GET("", s.handlerRegistry.settingHandler.ListAccessToken)
		accessTokenGroup.POST("", s.handlerRegistry.settingHandler.CreateAccessToken)
		accessTokenGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateAccessToken)
		accessTokenGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateAccessTokenMeta)
		accessTokenGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteAccessToken)
		// Test connection
		accessTokenGroup.POST("/test-conn", s.handlerRegistry.settingHandler.TestAccessTokenConn)
	}

	{ // aws group
		awsGroup := settingGroup.Group("/aws")
		awsGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetAWS)
		awsGroup.GET("", s.handlerRegistry.settingHandler.ListAWS)
		awsGroup.POST("", s.handlerRegistry.settingHandler.CreateAWS)
		awsGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateAWS)
		awsGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateAWSMeta)
		awsGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteAWS)
	}

	{ // aws s3 group
		awsS3Group := settingGroup.Group("/aws-s3")
		awsS3Group.GET("/:settingID", s.handlerRegistry.settingHandler.GetAWSS3)
		awsS3Group.GET("", s.handlerRegistry.settingHandler.ListAWSS3)
		awsS3Group.POST("", s.handlerRegistry.settingHandler.CreateAWSS3)
		awsS3Group.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateAWSS3)
		awsS3Group.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateAWSS3Meta)
		awsS3Group.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteAWSS3)
		// Test connection
		awsS3Group.POST("/test-conn", s.handlerRegistry.settingHandler.TestAWSS3Conn)
	}

	{ // ssh key group
		sshKeyGroup := settingGroup.Group("/ssh-keys")
		sshKeyGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetSSHKey)
		sshKeyGroup.GET("", s.handlerRegistry.settingHandler.ListSSHKey)
		sshKeyGroup.POST("", s.handlerRegistry.settingHandler.CreateSSHKey)
		sshKeyGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateSSHKey)
		sshKeyGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateSSHKeyMeta)
		sshKeyGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteSSHKey)
	}

	{ // IM service group
		imServiceGroup := settingGroup.Group("/im-services")
		imServiceGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetIMService)
		imServiceGroup.GET("", s.handlerRegistry.settingHandler.ListIMService)
		imServiceGroup.POST("", s.handlerRegistry.settingHandler.CreateIMService)
		imServiceGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateIMService)
		imServiceGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateIMServiceMeta)
		imServiceGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteIMService)
		// Test connection
		imServiceGroup.POST("/test-send-msg", s.handlerRegistry.settingHandler.TestSendInstantMsg)
	}

	{ // registry auth group
		registryAuthGroup := settingGroup.Group("/registry-auth")
		registryAuthGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetRegistryAuth)
		registryAuthGroup.GET("", s.handlerRegistry.settingHandler.ListRegistryAuth)
		registryAuthGroup.POST("", s.handlerRegistry.settingHandler.CreateRegistryAuth)
		registryAuthGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateRegistryAuth)
		registryAuthGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateRegistryAuthMeta)
		registryAuthGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteRegistryAuth)
		// Test connection
		registryAuthGroup.POST("/test-conn", s.handlerRegistry.settingHandler.TestRegistryAuthConn)
	}

	{ // basic auth group
		basicAuthGroup := settingGroup.Group("/basic-auth")
		basicAuthGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetBasicAuth)
		basicAuthGroup.GET("", s.handlerRegistry.settingHandler.ListBasicAuth)
		basicAuthGroup.POST("", s.handlerRegistry.settingHandler.CreateBasicAuth)
		basicAuthGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateBasicAuth)
		basicAuthGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateBasicAuthMeta)
		basicAuthGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteBasicAuth)
	}

	{ // ssl group
		sslGroup := settingGroup.Group("/ssls")
		sslGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetSSL)
		sslGroup.GET("", s.handlerRegistry.settingHandler.ListSSL)
		sslGroup.POST("", s.handlerRegistry.settingHandler.CreateSSL)
		sslGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateSSL)
		sslGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateSSLMeta)
		sslGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteSSL)
	}

	{ // email group
		emailGroup := settingGroup.Group("/emails")
		emailGroup.GET("/:settingID", s.handlerRegistry.settingHandler.GetEmail)
		emailGroup.GET("", s.handlerRegistry.settingHandler.ListEmail)
		emailGroup.POST("", s.handlerRegistry.settingHandler.CreateEmail)
		emailGroup.PUT("/:settingID", s.handlerRegistry.settingHandler.UpdateEmail)
		emailGroup.PUT("/:settingID/meta", s.handlerRegistry.settingHandler.UpdateEmailMeta)
		emailGroup.DELETE("/:settingID", s.handlerRegistry.settingHandler.DeleteEmail)
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
