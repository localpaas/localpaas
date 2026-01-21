package basesettinghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/discorduc/discorddto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc/gittokendto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc/slackdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (h *BaseSettingHandler) DeleteSetting(
	ctx *gin.Context,
	resType base.ResourceType,
	scope base.SettingScope,
) {
	var auth *basedto.Auth
	var projectID, appID, itemID string
	var err error

	switch scope {
	case base.SettingScopeGlobal:
		auth, itemID, err = h.GetAuthGlobalSettings(ctx, resType, base.ActionTypeDelete, true)
	case base.SettingScopeProject:
		auth, projectID, itemID, err = h.GetAuthProjectSettings(ctx, base.ActionTypeWrite, true)
	case base.SettingScopeApp:
		auth, projectID, appID, itemID, err = h.GetAuthAppSettings(ctx, base.ActionTypeWrite, true)
	case base.SettingScopeUser:
		auth, itemID, err = h.GetAuthUserSettings(ctx, base.ActionTypeWrite, true)
	}
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	var req any
	var ucFunc func() (any, error)
	reqCtx := h.RequestCtx(ctx)

	switch resType { //nolint:exhaustive
	case base.ResourceTypeBasicAuth:
		r := basicauthdto.NewDeleteBasicAuthReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.BasicAuthUC.DeleteBasicAuth(reqCtx, auth, r) }

	case base.ResourceTypeGithubApp:
		r := githubappdto.NewDeleteGithubAppReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.GithubAppUC.DeleteGithubApp(reqCtx, auth, r) }

	case base.ResourceTypeGitToken:
		r := gittokendto.NewDeleteGitTokenReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.GitTokenUC.DeleteGitToken(reqCtx, auth, r) }

	case base.ResourceTypeOAuth:
		r := oauthdto.NewDeleteOAuthReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.OAuthUC.DeleteOAuth(reqCtx, auth, r) }

	case base.ResourceTypeRegistryAuth:
		r := registryauthdto.NewDeleteRegistryAuthReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.RegistryAuthUC.DeleteRegistryAuth(reqCtx, auth, r) }

	case base.ResourceTypeS3Storage:
		r := s3storagedto.NewDeleteS3StorageReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.S3StorageUC.DeleteS3Storage(reqCtx, auth, r) }

	case base.ResourceTypeSSHKey:
		r := sshkeydto.NewDeleteSSHKeyReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SSHKeyUC.DeleteSSHKey(reqCtx, auth, r) }

	case base.ResourceTypeSSL:
		r := ssldto.NewDeleteSslReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SSLUC.DeleteSsl(reqCtx, auth, r) }

	case base.ResourceTypeCronJob:
		r := cronjobdto.NewDeleteCronJobReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.CronJobUC.DeleteCronJob(reqCtx, auth, r) }

	case base.ResourceTypeSecret:
		r := secretdto.NewDeleteSecretReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SecretUC.DeleteSecret(reqCtx, auth, r) }

	case base.ResourceTypeAPIKey:
		r := apikeydto.NewDeleteAPIKeyReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.APIKeyUC.DeleteAPIKey(reqCtx, auth, r) }

	case base.ResourceTypeSlack:
		r := slackdto.NewDeleteSlackReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SlackUC.DeleteSlack(reqCtx, auth, r) }

	case base.ResourceTypeDiscord:
		r := discorddto.NewDeleteDiscordReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.DiscordUC.DeleteDiscord(reqCtx, auth, r) }
	}

	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := ucFunc()
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
