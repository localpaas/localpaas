package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
)

// ListGithubApp Lists github-app providers
// @Summary Lists github-app providers
// @Description Lists github-app providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderGithubApp
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} githubappdto.ListGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/github-apps [get]
func (h *ProvidersHandler) ListGithubApp(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// GetGithubApp Gets github-app provider details
// @Summary Gets github-app provider details
// @Description Gets github-app provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderGithubApp
// @Param   id path string true "provider ID"
// @Success 200 {object} githubappdto.GetGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/github-apps/{id} [get]
func (h *ProvidersHandler) GetGithubApp(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// CreateGithubApp Creates a new github-app provider
// @Summary Creates a new github-app provider
// @Description Creates a new github-app provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderGithubApp
// @Param   body body githubappdto.CreateGithubAppReq true "request data"
// @Success 201 {object} githubappdto.CreateGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/github-apps [post]
func (h *ProvidersHandler) CreateGithubApp(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// UpdateGithubApp Updates github-app
// @Summary Updates github-app
// @Description Updates github-app
// @Tags    global_providers
// @Produce json
// @Id      updateProviderGithubApp
// @Param   id path string true "provider ID"
// @Param   body body githubappdto.UpdateGithubAppReq true "request data"
// @Success 200 {object} githubappdto.UpdateGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/github-apps/{id} [put]
func (h *ProvidersHandler) UpdateGithubApp(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// UpdateGithubAppMeta Updates github-app meta
// @Summary Updates github-app meta
// @Description Updates github-app meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderGithubAppMeta
// @Param   id path string true "provider ID"
// @Param   body body githubappdto.UpdateGithubAppMetaReq true "request data"
// @Success 200 {object} githubappdto.UpdateGithubAppMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/github-apps/{id}/meta [put]
func (h *ProvidersHandler) UpdateGithubAppMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// DeleteGithubApp Deletes github-app provider
// @Summary Deletes github-app provider
// @Description Deletes github-app provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderGithubApp
// @Param   id path string true "provider ID"
// @Success 200 {object} githubappdto.DeleteGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/github-apps/{id} [delete]
func (h *ProvidersHandler) DeleteGithubApp(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeGlobal)
}

// TestGithubAppConn Test github app connection
// @Summary Test github app connection
// @Description Test github app connection
// @Tags    global_providers
// @Produce json
// @Id      testGithubAppConn
// @Param   body body githubappdto.TestGithubAppConnReq true "request data"
// @Success 200 {object} githubappdto.TestGithubAppConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/github-apps/test-conn [post]
func (h *ProvidersHandler) TestGithubAppConn(ctx *gin.Context) {
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
// @Tags    global_providers
// @Produce json
// @Id      listAppInstallation
// @Param   body body githubappdto.ListAppInstallationReq true "request data"
// @Success 200 {object} githubappdto.ListAppInstallationResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/github-apps/installations/list [post]
func (h *ProvidersHandler) ListAppInstallation(ctx *gin.Context) {
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
