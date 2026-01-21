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

type ListSettingOptions struct {
	PreRequestHandler func(auth *basedto.Auth, req any) error
}

type ListSettingOption func(*ListSettingOptions)

func ListSettingPreRequestHandler(fn func(auth *basedto.Auth, req any) error) ListSettingOption {
	return func(opts *ListSettingOptions) {
		opts.PreRequestHandler = fn
	}
}

func (h *BaseSettingHandler) ListSetting(
	ctx *gin.Context,
	resType base.ResourceType,
	scope base.SettingScope,
	opts ...ListSettingOption,
) {
	var auth *basedto.Auth
	var projectID, appID string
	var err error

	options := &ListSettingOptions{}
	for _, o := range opts {
		o(options)
	}

	switch scope {
	case base.SettingScopeGlobal:
		auth, _, err = h.GetAuthGlobalSettings(ctx, resType, base.ActionTypeRead, "")
	case base.SettingScopeProject:
		auth, projectID, _, err = h.GetAuthProjectSettings(ctx, base.ActionTypeRead, "")
	case base.SettingScopeApp:
		auth, projectID, appID, _, err = h.GetAuthAppSettings(ctx, base.ActionTypeRead, "")
	case base.SettingScopeUser:
		auth, _, err = h.GetAuthUserSettings(ctx, base.ActionTypeRead, "")
	}
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	var req any
	var paging *basedto.Paging
	var ucFunc func() (any, error)
	reqCtx := h.RequestCtx(ctx)

	switch resType { //nolint:exhaustive
	case base.ResourceTypeBasicAuth:
		r := basicauthdto.NewListBasicAuthReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.BasicAuthUC.ListBasicAuth(reqCtx, auth, r) }

	case base.ResourceTypeGithubApp:
		r := githubappdto.NewListGithubAppReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.GithubAppUC.ListGithubApp(reqCtx, auth, r) }

	case base.ResourceTypeGitToken:
		r := gittokendto.NewListGitTokenReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.GitTokenUC.ListGitToken(reqCtx, auth, r) }

	case base.ResourceTypeOAuth:
		r := oauthdto.NewListOAuthReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.OAuthUC.ListOAuth(reqCtx, auth, r) }

	case base.ResourceTypeRegistryAuth:
		r := registryauthdto.NewListRegistryAuthReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.RegistryAuthUC.ListRegistryAuth(reqCtx, auth, r) }

	case base.ResourceTypeS3Storage:
		r := s3storagedto.NewListS3StorageReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.S3StorageUC.ListS3Storage(reqCtx, auth, r) }

	case base.ResourceTypeSSHKey:
		r := sshkeydto.NewListSSHKeyReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SSHKeyUC.ListSSHKey(reqCtx, auth, r) }

	case base.ResourceTypeSSL:
		r := ssldto.NewListSslReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SSLUC.ListSsl(reqCtx, auth, r) }

	case base.ResourceTypeCronJob:
		r := cronjobdto.NewListCronJobReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.CronJobUC.ListCronJob(reqCtx, auth, r) }

	case base.ResourceTypeSecret:
		r := secretdto.NewListSecretReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SecretUC.ListSecret(reqCtx, auth, r) }

	case base.ResourceTypeAPIKey:
		r := apikeydto.NewListAPIKeyReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.APIKeyUC.ListAPIKey(reqCtx, auth, r) }

	case base.ResourceTypeSlack:
		r := slackdto.NewListSlackReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SlackUC.ListSlack(reqCtx, auth, r) }

	case base.ResourceTypeDiscord:
		r := discorddto.NewListDiscordReq()
		r.Scope, r.ProjectID, r.AppID = scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.DiscordUC.ListDiscord(reqCtx, auth, r) }
	}

	if err = h.ParseAndValidateRequest(ctx, req, paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	if options.PreRequestHandler != nil {
		if err = options.PreRequestHandler(auth, req); err != nil {
			h.RenderError(ctx, err)
			return
		}
	}

	resp, err := ucFunc()
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
