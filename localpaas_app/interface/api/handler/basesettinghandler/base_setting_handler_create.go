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

func (h *BaseSettingHandler) CreateSetting(
	ctx *gin.Context,
	resType base.ResourceType,
	scope base.SettingScope,
) {
	var auth *basedto.Auth
	var projectID, appID string
	var err error

	switch scope {
	case base.SettingScopeGlobal:
		auth, _, err = h.GetAuthGlobalSettings(ctx, resType, base.ActionTypeWrite, false)
	case base.SettingScopeProject:
		auth, projectID, _, err = h.GetAuthProjectSettings(ctx, base.ActionTypeWrite, false)
	case base.SettingScopeApp:
		auth, projectID, appID, _, err = h.GetAuthAppSettings(ctx, base.ActionTypeWrite, false)
	case base.SettingScopeUser:
		auth, _, err = h.GetAuthUserSettings(ctx, base.ActionTypeWrite, false)
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
		r := basicauthdto.NewCreateBasicAuthReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.BasicAuthUC.CreateBasicAuth(reqCtx, auth, r) }

	case base.ResourceTypeGithubApp:
		r := githubappdto.NewCreateGithubAppReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.GithubAppUC.CreateGithubApp(reqCtx, auth, r) }

	case base.ResourceTypeGitToken:
		r := gittokendto.NewCreateGitTokenReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.GitTokenUC.CreateGitToken(reqCtx, auth, r) }

	case base.ResourceTypeOAuth:
		r := oauthdto.NewCreateOAuthReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.OAuthUC.CreateOAuth(reqCtx, auth, r) }

	case base.ResourceTypeRegistryAuth:
		r := registryauthdto.NewCreateRegistryAuthReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.RegistryAuthUC.CreateRegistryAuth(reqCtx, auth, r) }

	case base.ResourceTypeS3Storage:
		r := s3storagedto.NewCreateS3StorageReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.S3StorageUC.CreateS3Storage(reqCtx, auth, r) }

	case base.ResourceTypeSSHKey:
		r := sshkeydto.NewCreateSSHKeyReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SSHKeyUC.CreateSSHKey(reqCtx, auth, r) }

	case base.ResourceTypeSSL:
		r := ssldto.NewCreateSslReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SSLUC.CreateSsl(reqCtx, auth, r) }

	case base.ResourceTypeCronJob:
		r := cronjobdto.NewCreateCronJobReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.CronJobUC.CreateCronJob(reqCtx, auth, r) }

	case base.ResourceTypeSecret:
		r := secretdto.NewCreateSecretReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SecretUC.CreateSecret(reqCtx, auth, r) }

	case base.ResourceTypeAPIKey:
		r := apikeydto.NewCreateAPIKeyReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.APIKeyUC.CreateAPIKey(reqCtx, auth, r) }

	case base.ResourceTypeSlack:
		r := slackdto.NewCreateSlackReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SlackUC.CreateSlack(reqCtx, auth, r) }

	case base.ResourceTypeDiscord:
		r := discorddto.NewCreateDiscordReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.DiscordUC.CreateDiscord(reqCtx, auth, r) }
	}

	if err = h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := ucFunc()
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}
