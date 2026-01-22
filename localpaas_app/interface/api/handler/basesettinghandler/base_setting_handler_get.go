package basesettinghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc/slackdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

type GetSettingOptions struct {
	PreRequestHandler func(auth *basedto.Auth, req any) error
}

type GetSettingOption func(*GetSettingOptions)

func GetSettingPreRequestHandler(fn func(auth *basedto.Auth, req any) error) GetSettingOption {
	return func(opts *GetSettingOptions) {
		opts.PreRequestHandler = fn
	}
}

func (h *BaseSettingHandler) GetSetting(
	ctx *gin.Context,
	resType base.ResourceType,
	scope base.SettingScope,
	opts ...GetSettingOption,
) {
	var auth *basedto.Auth
	var objectID, parentObjectID, itemID string
	var err error

	options := &GetSettingOptions{}
	for _, o := range opts {
		o(options)
	}

	switch scope {
	case base.SettingScopeGlobal:
		auth, itemID, err = h.GetAuthGlobalSettings(ctx, resType, base.ActionTypeRead, "id")
	case base.SettingScopeProject:
		auth, objectID, itemID, err = h.GetAuthProjectSettings(ctx, base.ActionTypeRead, "id")
	case base.SettingScopeApp:
		auth, parentObjectID, objectID, itemID, err = h.GetAuthAppSettings(ctx, base.ActionTypeRead, "id")
	case base.SettingScopeUser:
		auth, objectID, itemID, err = h.GetAuthUserSettings(ctx, base.ActionTypeRead, "id")
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
		r := basicauthdto.NewGetBasicAuthReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.BasicAuthUC.GetBasicAuth(reqCtx, auth, r) }

	case base.ResourceTypeGithubApp:
		r := githubappdto.NewGetGithubAppReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.GithubAppUC.GetGithubApp(reqCtx, auth, r) }

	case base.ResourceTypeGitToken:
		r := gittokendto.NewGetGitTokenReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.GitTokenUC.GetGitToken(reqCtx, auth, r) }

	case base.ResourceTypeOAuth:
		r := oauthdto.NewGetOAuthReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.OAuthUC.GetOAuth(reqCtx, auth, r) }

	case base.ResourceTypeRegistryAuth:
		r := registryauthdto.NewGetRegistryAuthReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.RegistryAuthUC.GetRegistryAuth(reqCtx, auth, r) }

	case base.ResourceTypeS3Storage:
		r := s3storagedto.NewGetS3StorageReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.S3StorageUC.GetS3Storage(reqCtx, auth, r) }

	case base.ResourceTypeSSHKey:
		r := sshkeydto.NewGetSSHKeyReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SSHKeyUC.GetSSHKey(reqCtx, auth, r) }

	case base.ResourceTypeSSL:
		r := ssldto.NewGetSslReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SSLUC.GetSsl(reqCtx, auth, r) }

	case base.ResourceTypeCronJob:
		r := cronjobdto.NewGetCronJobReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.CronJobUC.GetCronJob(reqCtx, auth, r) }

	case base.ResourceTypeSecret:
		// NOTE: not implemented
		err = apperrors.NewNotImplementedNT()

	case base.ResourceTypeAPIKey:
		r := apikeydto.NewGetAPIKeyReq()
		r.ID = itemID
		req, ucFunc = r, func() (any, error) { return h.APIKeyUC.GetAPIKey(reqCtx, auth, r) }

	case base.ResourceTypeSlack:
		r := slackdto.NewGetSlackReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SlackUC.GetSlack(reqCtx, auth, r) }

	case base.ResourceTypeDiscord:
		r := discorddto.NewGetDiscordReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.DiscordUC.GetDiscord(reqCtx, auth, r) }
	}
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
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
