package basesettinghandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc"
)

type BaseSettingHandler struct {
	*handler.BaseHandler
	AuthHandler    *authhandler.AuthHandler
	OAuthUC        *oauthuc.OAuthUC
	S3StorageUC    *s3storageuc.S3StorageUC
	SSHKeyUC       *sshkeyuc.SSHKeyUC
	IMServiceUC    *imserviceuc.IMServiceUC
	RegistryAuthUC *registryauthuc.RegistryAuthUC
	BasicAuthUC    *basicauthuc.BasicAuthUC
	SSLUC          *ssluc.SslUC
	GithubAppUC    *githubappuc.GithubAppUC
	GitTokenUC     *gittokenuc.GitTokenUC
	CronJobUC      *cronjobuc.CronJobUC
	SecretUC       *secretuc.SecretUC
	EmailUC        *emailuc.EmailUC
	APIKeyUC       *apikeyuc.APIKeyUC
}

func NewBaseSettingHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	oauthUC *oauthuc.OAuthUC,
	s3StorageUC *s3storageuc.S3StorageUC,
	sshKeyUC *sshkeyuc.SSHKeyUC,
	imServiceUC *imserviceuc.IMServiceUC,
	registryAuthUC *registryauthuc.RegistryAuthUC,
	basicAuthUC *basicauthuc.BasicAuthUC,
	sslUC *ssluc.SslUC,
	githubAppUC *githubappuc.GithubAppUC,
	gitTokenUC *gittokenuc.GitTokenUC,
	cronJobUC *cronjobuc.CronJobUC,
	secretUC *secretuc.SecretUC,
	emailUC *emailuc.EmailUC,
	apiKeyUC *apikeyuc.APIKeyUC,
) *BaseSettingHandler {
	return &BaseSettingHandler{
		BaseHandler:    baseHandler,
		AuthHandler:    authHandler,
		OAuthUC:        oauthUC,
		S3StorageUC:    s3StorageUC,
		SSHKeyUC:       sshKeyUC,
		IMServiceUC:    imServiceUC,
		RegistryAuthUC: registryAuthUC,
		BasicAuthUC:    basicAuthUC,
		SSLUC:          sslUC,
		GithubAppUC:    githubAppUC,
		GitTokenUC:     gitTokenUC,
		CronJobUC:      cronJobUC,
		SecretUC:       secretUC,
		EmailUC:        emailUC,
		APIKeyUC:       apiKeyUC,
	}
}
