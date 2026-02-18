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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/notificationuc/notificationdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/repowebhookuc/repowebhookdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

//nolint:funlen
func (h *BaseSettingHandler) UpdateSettingMeta(
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
	case base.SettingScopeNone:
		err = apperrors.NewUnsupported().WithMsgLog("Setting scope 'none' is not supported")
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
		r := basicauthdto.NewUpdateBasicAuthMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.BasicAuthUC.UpdateBasicAuthMeta(reqCtx, auth, r) }

	case base.ResourceTypeGithubApp:
		r := githubappdto.NewUpdateGithubAppMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.GithubAppUC.UpdateGithubAppMeta(reqCtx, auth, r) }

	case base.ResourceTypeAccessToken:
		r := accesstokendto.NewUpdateAccessTokenMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AccessTokenUC.UpdateAccessTokenMeta(reqCtx, auth, r) }

	case base.ResourceTypeOAuth:
		r := oauthdto.NewUpdateOAuthMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.OAuthUC.UpdateOAuthMeta(reqCtx, auth, r) }

	case base.ResourceTypeRegistryAuth:
		r := registryauthdto.NewUpdateRegistryAuthMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.RegistryAuthUC.UpdateRegistryAuthMeta(reqCtx, auth, r) }

	case base.ResourceTypeAWS:
		r := awsdto.NewUpdateAWSMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AWSUC.UpdateAWSMeta(reqCtx, auth, r) }

	case base.ResourceTypeAWSS3:
		r := awss3dto.NewUpdateAWSS3MetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.AWSS3UC.UpdateAWSS3Meta(reqCtx, auth, r) }

	case base.ResourceTypeSSHKey:
		r := sshkeydto.NewUpdateSSHKeyMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SSHKeyUC.UpdateSSHKeyMeta(reqCtx, auth, r) }

	case base.ResourceTypeSSL:
		r := ssldto.NewUpdateSSLMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SSLUC.UpdateSSLMeta(reqCtx, auth, r) }

	case base.ResourceTypeCronJob:
		r := cronjobdto.NewUpdateCronJobMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.CronJobUC.UpdateCronJobMeta(reqCtx, auth, r) }

	case base.ResourceTypeHealthcheck:
		r := healthcheckdto.NewUpdateHealthcheckMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.HealthcheckUC.UpdateHealthcheckMeta(reqCtx, auth, r) }

	case base.ResourceTypeSecret:
		r := secretdto.NewUpdateSecretMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.SecretUC.UpdateSecretMeta(reqCtx, auth, r) }

	case base.ResourceTypeAPIKey:
		r := apikeydto.NewUpdateAPIKeyMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.APIKeyUC.UpdateAPIKeyMeta(reqCtx, auth, r) }

	case base.ResourceTypeIMService:
		r := imservicedto.NewUpdateIMServiceMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.IMServiceUC.UpdateIMServiceMeta(reqCtx, auth, r) }

	case base.ResourceTypeEmail:
		r := emaildto.NewUpdateEmailMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.EmailUC.UpdateEmailMeta(reqCtx, auth, r) }

	case base.ResourceTypeRepoWebhook:
		r := repowebhookdto.NewUpdateRepoWebhookMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.RepoWebhookUC.UpdateRepoWebhookMeta(reqCtx, auth, r) }

	case base.ResourceTypeNotification:
		r := notificationdto.NewUpdateNotificationMetaReq()
		r.ID, r.Scope, r.ObjectID, r.ParentObjectID = itemID, scope, objectID, parentObjectID
		req, ucFunc = r, func() (any, error) { return h.NotificationUC.UpdateNotificationMeta(reqCtx, auth, r) }
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
