package basesettinghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/accesstokenuc/accesstokendto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awss3uc/awss3dto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awsuc/awsdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/healthcheckuc/healthcheckdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc/imservicedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/repowebhookuc/repowebhookdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

type UpdateSettingOptions struct {
	PreRequestHandler func(auth *basedto.Auth, req any) error
}

type UpdateSettingOption func(*UpdateSettingOptions)

func UpdateSettingPreRequestHandler(fn func(auth *basedto.Auth, req any) error) UpdateSettingOption {
	return func(opts *UpdateSettingOptions) {
		opts.PreRequestHandler = fn
	}
}

//nolint:funlen
func (h *BaseSettingHandler) UpdateSetting(
	ctx *gin.Context,
	resType base.ResourceType,
	scope base.SettingScope,
	opts ...UpdateSettingOption,
) {
	var auth *basedto.Auth
	var objectID, parentObjectID, itemID string
	var err error

	options := &UpdateSettingOptions{}
	for _, o := range opts {
		o(options)
	}

	switch scope {
	case base.SettingScopeGlobal:
		auth, itemID, err = h.GetAuthGlobalSettings(ctx, resType, base.ActionTypeWrite, "itemID")
	case base.SettingScopeProject:
		auth, objectID, itemID, err = h.GetAuthProjectSettings(ctx, base.ActionTypeWrite, "itemID")
	case base.SettingScopeApp:
		auth, parentObjectID, objectID, itemID, err = h.GetAuthAppSettings(ctx, base.ActionTypeWrite, "itemID")
	case base.SettingScopeUser:
		auth, objectID, itemID, err = h.GetAuthUserSettings(ctx, base.ActionTypeWrite, "itemID")
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
		r := basicauthdto.NewUpdateBasicAuthReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.BasicAuthUC.UpdateBasicAuth(reqCtx, auth, r) }

	case base.ResourceTypeGithubApp:
		r := githubappdto.NewUpdateGithubAppReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.GithubAppUC.UpdateGithubApp(reqCtx, auth, r) }

	case base.ResourceTypeAccessToken:
		r := accesstokendto.NewUpdateAccessTokenReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AccessTokenUC.UpdateAccessToken(reqCtx, auth, r) }

	case base.ResourceTypeOAuth:
		r := oauthdto.NewUpdateOAuthReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.OAuthUC.UpdateOAuth(reqCtx, auth, r) }

	case base.ResourceTypeRegistryAuth:
		r := registryauthdto.NewUpdateRegistryAuthReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.RegistryAuthUC.UpdateRegistryAuth(reqCtx, auth, r) }

	case base.ResourceTypeAWS:
		r := awsdto.NewUpdateAWSReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AWSUC.UpdateAWS(reqCtx, auth, r) }

	case base.ResourceTypeAWSS3:
		r := awss3dto.NewUpdateAWSS3Req()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AWSS3UC.UpdateAWSS3(reqCtx, auth, r) }

	case base.ResourceTypeSSHKey:
		r := sshkeydto.NewUpdateSSHKeyReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SSHKeyUC.UpdateSSHKey(reqCtx, auth, r) }

	case base.ResourceTypeSSL:
		r := ssldto.NewUpdateSSLReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SSLUC.UpdateSSL(reqCtx, auth, r) }

	case base.ResourceTypeCronJob:
		r := cronjobdto.NewUpdateCronJobReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.CronJobUC.UpdateCronJob(reqCtx, auth, r) }

	case base.ResourceTypeHealthcheck:
		r := healthcheckdto.NewUpdateHealthcheckReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.HealthcheckUC.UpdateHealthcheck(reqCtx, auth, r) }

	case base.ResourceTypeSecret:
		r := secretdto.NewUpdateSecretReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SecretUC.UpdateSecret(reqCtx, auth, r) }

	case base.ResourceTypeAPIKey:
		// NOTE: not implemented
		err = apperrors.NewNotImplementedNT()

	case base.ResourceTypeIMService:
		r := imservicedto.NewUpdateIMServiceReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.IMServiceUC.UpdateIMService(reqCtx, auth, r) }

	case base.ResourceTypeEmail:
		r := emaildto.NewUpdateEmailReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.EmailUC.UpdateEmail(reqCtx, auth, r) }

	case base.ResourceTypeRepoWebhook:
		r := repowebhookdto.NewUpdateRepoWebhookReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.RepoWebhookUC.UpdateRepoWebhook(reqCtx, auth, r) }
	}
	if err != nil {
		h.RenderError(ctx, err)
		return
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

	ctx.JSON(http.StatusOK, resp)
}
