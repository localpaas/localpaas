package providershandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/s3storageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/secretuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/slackuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/sshkeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc"
)

type ProvidersHandler struct {
	*handler.BaseHandler
	authHandler    *authhandler.AuthHandler
	oauthUC        *oauthuc.OAuthUC
	s3StorageUC    *s3storageuc.S3StorageUC
	sshKeyUC       *sshkeyuc.SSHKeyUC
	secretUC       *secretuc.SecretUC
	slackUC        *slackuc.SlackUC
	discordUC      *discorduc.DiscordUC
	registryAuthUC *registryauthuc.RegistryAuthUC
	basicAuthUC    *basicauthuc.BasicAuthUC
	sslUC          *ssluc.SslUC
	githubAppUC    *githubappuc.GithubAppUC
	gitTokenUC     *gittokenuc.GitTokenUC
	cronJobUC      *cronjobuc.CronJobUC
}

func NewProvidersHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	oauthUC *oauthuc.OAuthUC,
	s3StorageUC *s3storageuc.S3StorageUC,
	sshKeyUC *sshkeyuc.SSHKeyUC,
	secretUC *secretuc.SecretUC,
	slackUC *slackuc.SlackUC,
	discordUC *discorduc.DiscordUC,
	registryAuthUC *registryauthuc.RegistryAuthUC,
	basicAuthUC *basicauthuc.BasicAuthUC,
	sslUC *ssluc.SslUC,
	githubAppUC *githubappuc.GithubAppUC,
	gitTokenUC *gittokenuc.GitTokenUC,
	cronJobUC *cronjobuc.CronJobUC,
) *ProvidersHandler {
	return &ProvidersHandler{
		BaseHandler:    baseHandler,
		authHandler:    authHandler,
		oauthUC:        oauthUC,
		s3StorageUC:    s3StorageUC,
		sshKeyUC:       sshKeyUC,
		secretUC:       secretUC,
		slackUC:        slackUC,
		discordUC:      discordUC,
		registryAuthUC: registryAuthUC,
		basicAuthUC:    basicAuthUC,
		sslUC:          sslUC,
		githubAppUC:    githubAppUC,
		gitTokenUC:     gitTokenUC,
		cronJobUC:      cronJobUC,
	}
}
