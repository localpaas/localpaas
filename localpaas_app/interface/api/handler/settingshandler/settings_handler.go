package settingshandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc"
)

type SettingsHandler struct {
	*handler.BaseHandler
	authHandler    *authhandler.AuthHandler
	oauthUC        *oauthuc.OAuthUC
	s3StorageUC    *s3storageuc.S3StorageUC
	sshKeyUC       *sshkeyuc.SSHKeyUC
	secretUC       *secretuc.SecretUC
	slackUC        *slackuc.SlackUC
	registryAuthUC *registryauthuc.RegistryAuthUC
}

func NewSettingsHandler(
	authHandler *authhandler.AuthHandler,
	oauthUC *oauthuc.OAuthUC,
	s3StorageUC *s3storageuc.S3StorageUC,
	sshKeyUC *sshkeyuc.SSHKeyUC,
	secretUC *secretuc.SecretUC,
	slackUC *slackuc.SlackUC,
	registryAuthUC *registryauthuc.RegistryAuthUC,
) *SettingsHandler {
	hdl := &SettingsHandler{
		authHandler:    authHandler,
		oauthUC:        oauthUC,
		s3StorageUC:    s3StorageUC,
		sshKeyUC:       sshKeyUC,
		secretUC:       secretUC,
		slackUC:        slackUC,
		registryAuthUC: registryAuthUC,
	}
	return hdl
}
