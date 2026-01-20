package projecthandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/s3storageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/secretuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/slackuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/sshkeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc"
)

type ProjectHandler struct {
	*handler.BaseHandler
	authHandler    *authhandler.AuthHandler
	projectUC      *projectuc.ProjectUC
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

func NewProjectHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	projectUC *projectuc.ProjectUC,
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
) *ProjectHandler {
	return &ProjectHandler{
		BaseHandler:    baseHandler,
		authHandler:    authHandler,
		projectUC:      projectUC,
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
