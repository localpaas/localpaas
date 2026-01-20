package providershandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/discorduc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc"
)

type ProvidersHandler struct {
	*handler.BaseHandler
	authHandler    *authhandler.AuthHandler
	oauthUC        *oauthuc.OAuthUC
	s3StorageUC    *s3storageuc.S3StorageUC
	sshKeyUC       *sshkeyuc.SSHKeyUC
	slackUC        *slackuc.SlackUC
	discordUC      *discorduc.DiscordUC
	registryAuthUC *registryauthuc.RegistryAuthUC
	basicAuthUC    *basicauthuc.BasicAuthUC
	sslUC          *ssluc.SslUC
	githubAppUC    *githubappuc.GithubAppUC
	gitTokenUC     *gittokenuc.GitTokenUC
}

func NewProvidersHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	oauthUC *oauthuc.OAuthUC,
	s3StorageUC *s3storageuc.S3StorageUC,
	sshKeyUC *sshkeyuc.SSHKeyUC,
	slackUC *slackuc.SlackUC,
	discordUC *discorduc.DiscordUC,
	registryAuthUC *registryauthuc.RegistryAuthUC,
	basicAuthUC *basicauthuc.BasicAuthUC,
	sslUC *ssluc.SslUC,
	githubAppUC *githubappuc.GithubAppUC,
	gitTokenUC *gittokenuc.GitTokenUC,
) *ProvidersHandler {
	return &ProvidersHandler{
		BaseHandler:    baseHandler,
		authHandler:    authHandler,
		oauthUC:        oauthUC,
		s3StorageUC:    s3StorageUC,
		sshKeyUC:       sshKeyUC,
		slackUC:        slackUC,
		discordUC:      discordUC,
		registryAuthUC: registryAuthUC,
		basicAuthUC:    basicAuthUC,
		sslUC:          sslUC,
		githubAppUC:    githubAppUC,
		gitTokenUC:     gitTokenUC,
	}
}
