package projecthandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/discorduc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc"
)

type ProjectHandler struct {
	*basesettinghandler.BaseSettingHandler
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
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
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
		BaseSettingHandler: baseSettingHandler,
		authHandler:        authHandler,
		projectUC:          projectUC,
		s3StorageUC:        s3StorageUC,
		sshKeyUC:           sshKeyUC,
		secretUC:           secretUC,
		slackUC:            slackUC,
		discordUC:          discordUC,
		registryAuthUC:     registryAuthUC,
		basicAuthUC:        basicAuthUC,
		sslUC:              sslUC,
		githubAppUC:        githubAppUC,
		gitTokenUC:         gitTokenUC,
		cronJobUC:          cronJobUC,
	}
}
