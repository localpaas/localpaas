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
)

func (h *BaseSettingHandler) UpdateSetting(ctx *gin.Context, resType base.ResourceType, scope base.SettingScope) {
	var auth *basedto.Auth
	var projectID, appID, itemID string
	var err error

	switch scope {
	case base.SettingScopeGlobal:
		auth, itemID, err = h.GetAuthGlobalSettings(ctx, resType, base.ActionTypeWrite, true)
	case base.SettingScopeUser:
		auth, itemID, err = h.GetAuthUserSettings(ctx, base.ActionTypeWrite, true)
	case base.SettingScopeProject:
		auth, projectID, itemID, err = h.GetAuthProjectSettings(ctx, base.ActionTypeWrite, true)
	case base.SettingScopeApp:
		auth, projectID, appID, itemID, err = h.GetAuthAppSettings(ctx, base.ActionTypeWrite, true)
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
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.BasicAuthUC.UpdateBasicAuth(reqCtx, auth, r) }

	case base.ResourceTypeGithubApp:
		r := githubappdto.NewUpdateGithubAppReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.GithubAppUC.UpdateGithubApp(reqCtx, auth, r) }

	case base.ResourceTypeGitToken:
		r := gittokendto.NewUpdateGitTokenReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.GitTokenUC.UpdateGitToken(reqCtx, auth, r) }

	case base.ResourceTypeOAuth:
		r := oauthdto.NewUpdateOAuthReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.OAuthUC.UpdateOAuth(reqCtx, auth, r) }

	case base.ResourceTypeRegistryAuth:
		r := registryauthdto.NewUpdateRegistryAuthReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.RegistryAuthUC.UpdateRegistryAuth(reqCtx, auth, r) }

	case base.ResourceTypeS3Storage:
		r := s3storagedto.NewUpdateS3StorageReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.S3StorageUC.UpdateS3Storage(reqCtx, auth, r) }

	case base.ResourceTypeSSHKey:
		r := sshkeydto.NewUpdateSSHKeyReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SSHKeyUC.UpdateSSHKey(reqCtx, auth, r) }

	case base.ResourceTypeSSL:
		r := ssldto.NewUpdateSslReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SSLUC.UpdateSsl(reqCtx, auth, r) }

	case base.ResourceTypeCronJob:
		r := cronjobdto.NewUpdateCronJobReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.CronJobUC.UpdateCronJob(reqCtx, auth, r) }

	case base.ResourceTypeSecret:
		// NOTE: not implemented
		err = apperrors.NewNotImplementedNT()

	case base.ResourceTypeAPIKey:
		// NOTE: not implemented
		err = apperrors.NewNotImplementedNT()

	case base.ResourceTypeSlack:
		r := slackdto.NewUpdateSlackReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.SlackUC.UpdateSlack(reqCtx, auth, r) }

	case base.ResourceTypeDiscord:
		r := discorddto.NewUpdateDiscordReq()
		r.ID, r.Scope, r.ProjectID, r.AppID = itemID, scope, projectID, appID
		req, ucFunc = r, func() (any, error) { return h.DiscordUC.UpdateDiscord(reqCtx, auth, r) }
	}
	if err != nil {
		h.RenderError(ctx, err)
		return
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

	ctx.JSON(http.StatusOK, resp)
}
