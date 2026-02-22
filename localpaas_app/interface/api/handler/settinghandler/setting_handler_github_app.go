package settinghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
)

const (
	htmlSuccess = "<html><body><h3>Success</h3></body></html>"
)

// ListGithubApp Lists github-app settings
// @Summary Lists github-app settings
// @Description Lists github-app settings
// @Tags    settings
// @Produce json
// @Id      listSettingGithubApp
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} githubappdto.ListGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps [get]
func (h *SettingHandler) ListGithubApp(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// GetGithubApp Gets github-app setting details
// @Summary Gets github-app setting details
// @Description Gets github-app setting details
// @Tags    settings
// @Produce json
// @Id      getSettingGithubApp
// @Param   itemID path string true "setting ID"
// @Success 200 {object} githubappdto.GetGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/{itemID} [get]
func (h *SettingHandler) GetGithubApp(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// CreateGithubApp Creates a new github-app setting
// @Summary Creates a new github-app setting
// @Description Creates a new github-app setting
// @Tags    settings
// @Produce json
// @Id      createSettingGithubApp
// @Param   body body githubappdto.CreateGithubAppReq true "request data"
// @Success 201 {object} githubappdto.CreateGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps [post]
func (h *SettingHandler) CreateGithubApp(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// UpdateGithubApp Updates github-app
// @Summary Updates github-app
// @Description Updates github-app
// @Tags    settings
// @Produce json
// @Id      updateSettingGithubApp
// @Param   itemID path string true "setting ID"
// @Param   body body githubappdto.UpdateGithubAppReq true "request data"
// @Success 200 {object} githubappdto.UpdateGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/{itemID} [put]
func (h *SettingHandler) UpdateGithubApp(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// UpdateGithubAppMeta Updates github-app meta
// @Summary Updates github-app meta
// @Description Updates github-app meta
// @Tags    settings
// @Produce json
// @Id      updateSettingGithubAppMeta
// @Param   itemID path string true "setting ID"
// @Param   body body githubappdto.UpdateGithubAppMetaReq true "request data"
// @Success 200 {object} githubappdto.UpdateGithubAppMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/{itemID}/meta [put]
func (h *SettingHandler) UpdateGithubAppMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// DeleteGithubApp Deletes github-app setting
// @Summary Deletes github-app setting
// @Description Deletes github-app setting
// @Tags    settings
// @Produce json
// @Id      deleteSettingGithubApp
// @Param   itemID path string true "setting ID"
// @Success 200 {object} githubappdto.DeleteGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/{itemID} [delete]
func (h *SettingHandler) DeleteGithubApp(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// TestGithubAppConn Test github app connection
// @Summary Test github app connection
// @Description Test github app connection
// @Tags    settings
// @Produce json
// @Id      testGithubAppConn
// @Param   body body githubappdto.TestGithubAppConnReq true "request data"
// @Success 200 {object} githubappdto.TestGithubAppConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/test-conn [post]
func (h *SettingHandler) TestGithubAppConn(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := githubappdto.NewTestGithubAppConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GithubAppUC.TestGithubAppConn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListAppInstallation List github app installation
// @Summary List github app installation
// @Description List github app installation
// @Tags    settings
// @Produce json
// @Id      listAppInstallation
// @Param   body body githubappdto.ListAppInstallationReq true "request data"
// @Success 200 {object} githubappdto.ListAppInstallationResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/installations/list [post]
func (h *SettingHandler) ListAppInstallation(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := githubappdto.NewListAppInstallationReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}
	if err = h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GithubAppUC.ListAppInstallation(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// BeginGithubAppManifestFlow Begins a github-app manifest flow
// @Summary Begins a github-app manifest flow
// @Description Begins a github-app manifest flow
// @Tags    settings
// @Produce json
// @Id      beginGithubAppManifestFlow
// @Param   body body githubappdto.BeginGithubAppManifestFlowReq true "request data"
// @Success 200 {object} githubappdto.BeginGithubAppManifestFlowResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/manifest-flow/begin [post]
func (h *SettingHandler) BeginGithubAppManifestFlow(ctx *gin.Context) {
	auth, _, err := h.GetAuthGlobalSettings(ctx, base.ResourceTypeGithubApp, base.ActionTypeWrite, "")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := githubappdto.NewBeginGithubAppManifestFlowReq()
	req.Scope = base.SettingScopeGlobal
	req.Type = base.SettingTypeGithubApp
	if err = h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GithubAppUC.BeginGithubAppManifestFlow(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// BeginGithubAppManifestFlowCreation Begins a github-app manifest flow of creation
// @Summary Begins a github-app manifest flow of creation
// @Description Begins a github-app manifest flow of creation
// @Tags    settings
// @Produce json
// @Id      beginGithubAppManifestFlowCreation
// @Param   itemID path string true "setting ID"
// @Success 200 "html page to redirect to github app creation page"
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/{itemID}/manifest-flow/begin [get]
func (h *SettingHandler) BeginGithubAppManifestFlowCreation(ctx *gin.Context) {
	itemID, err := h.ParseStringParam(ctx, "itemID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := githubappdto.NewBeginGithubAppManifestFlowCreationReq()
	req.SettingID = itemID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GithubAppUC.BeginGithubAppManifestFlowCreation(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.Data(http.StatusOK, "text/html", reflectutil.UnsafeStrToBytes(resp.Data.PageContent))
}

// SetupGithubAppManifestFlow Sets up a github-app manifest flow
// @Summary Sets up a github-app manifest flow
// @Description Sets up a github-app manifest flow
// @Tags    settings
// @Produce json
// @Id      setupGithubAppManifestFlow
// @Param   itemID path string true "setting ID"
// @Success 200 "html page to redirect to github app creation page"
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/github-apps/{itemID}/manifest-flow/setup [get]
func (h *SettingHandler) SetupGithubAppManifestFlow(ctx *gin.Context) {
	itemID, err := h.ParseStringParam(ctx, "itemID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := githubappdto.NewSetupGithubAppManifestFlowReq()
	req.SettingID = itemID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GithubAppUC.SetupGithubAppManifestFlow(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	if resp.Data.RedirectURL != "" {
		ctx.Redirect(http.StatusFound, resp.Data.RedirectURL)
	}

	ctx.Data(http.StatusOK, "text/html", []byte(htmlSuccess))
}
