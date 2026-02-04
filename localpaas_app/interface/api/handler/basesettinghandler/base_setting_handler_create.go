package basesettinghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/accesstokenuc/accesstokendto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awss3uc/awss3dto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awsuc/awsdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc/imservicedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/repowebhookuc/repowebhookdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

type CreateSettingOptions struct {
	PreRequestHandler func(auth *basedto.Auth, req any) error
}

type CreateSettingOption func(*CreateSettingOptions)

func CreateSettingPreRequestHandler(fn func(auth *basedto.Auth, req any) error) CreateSettingOption {
	return func(opts *CreateSettingOptions) {
		opts.PreRequestHandler = fn
	}
}

//nolint:funlen
func (h *BaseSettingHandler) CreateSetting(
	ctx *gin.Context,
	resType base.ResourceType,
	scope base.SettingScope,
	opts ...CreateSettingOption,
) {
	var auth *basedto.Auth
	var objectID, parentObjectID string
	var err error

	options := &CreateSettingOptions{}
	for _, o := range opts {
		o(options)
	}

	switch scope {
	case base.SettingScopeGlobal:
		auth, _, err = h.GetAuthGlobalSettings(ctx, resType, base.ActionTypeWrite, "")
	case base.SettingScopeProject:
		auth, objectID, _, err = h.GetAuthProjectSettings(ctx, base.ActionTypeWrite, "")
	case base.SettingScopeApp:
		auth, parentObjectID, objectID, _, err = h.GetAuthAppSettings(ctx, base.ActionTypeWrite, "")
	case base.SettingScopeUser:
		auth, objectID, _, err = h.GetAuthUserSettings(ctx, base.ActionTypeWrite, "")
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
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.BasicAuthUC.CreateBasicAuth(reqCtx, auth, r) }

	case base.ResourceTypeGithubApp:
		r := githubappdto.NewCreateGithubAppReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.GithubAppUC.CreateGithubApp(reqCtx, auth, r) }

	case base.ResourceTypeAccessToken:
		r := accesstokendto.NewCreateAccessTokenReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AccessTokenUC.CreateAccessToken(reqCtx, auth, r) }

	case base.ResourceTypeOAuth:
		r := oauthdto.NewCreateOAuthReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.OAuthUC.CreateOAuth(reqCtx, auth, r) }

	case base.ResourceTypeRegistryAuth:
		r := registryauthdto.NewCreateRegistryAuthReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.RegistryAuthUC.CreateRegistryAuth(reqCtx, auth, r) }

	case base.ResourceTypeAWS:
		r := awsdto.NewCreateAWSReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AWSUC.CreateAWS(reqCtx, auth, r) }

	case base.ResourceTypeAWSS3:
		r := awss3dto.NewCreateAWSS3Req()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AWSS3UC.CreateAWSS3(reqCtx, auth, r) }

	case base.ResourceTypeSSHKey:
		r := sshkeydto.NewCreateSSHKeyReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SSHKeyUC.CreateSSHKey(reqCtx, auth, r) }

	case base.ResourceTypeSSL:
		r := ssldto.NewCreateSSLReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SSLUC.CreateSSL(reqCtx, auth, r) }

	case base.ResourceTypeCronJob:
		r := cronjobdto.NewCreateCronJobReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.CronJobUC.CreateCronJob(reqCtx, auth, r) }

	case base.ResourceTypeSecret:
		r := secretdto.NewCreateSecretReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SecretUC.CreateSecret(reqCtx, auth, r) }

	case base.ResourceTypeAPIKey:
		r := apikeydto.NewCreateAPIKeyReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.APIKeyUC.CreateAPIKey(reqCtx, auth, r) }

	case base.ResourceTypeIMService:
		r := imservicedto.NewCreateIMServiceReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.IMServiceUC.CreateIMService(reqCtx, auth, r) }

	case base.ResourceTypeEmail:
		r := emaildto.NewCreateEmailReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.EmailUC.CreateEmail(reqCtx, auth, r) }

	case base.ResourceTypeRepoWebhook:
		r := repowebhookdto.NewCreateRepoWebhookReq()
		r.Scope, r.ObjectID, r.ParentObjectID = scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.RepoWebhookUC.CreateRepoWebhook(reqCtx, auth, r) }
	}

	if err = h.ParseAndValidateJSONBody(ctx, req); err != nil {
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

	ctx.JSON(http.StatusCreated, resp)
}
