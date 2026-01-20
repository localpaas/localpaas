package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc/gittokendto"
)

// ListGitToken Lists git-token providers
// @Summary Lists git-token providers
// @Description Lists git-token providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderGitToken
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} gittokendto.ListGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/git-tokens [get]
func (h *ProvidersHandler) ListGitToken(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeGitToken, base.SettingScopeGlobal)
}

// GetGitToken Gets git-token provider details
// @Summary Gets git-token provider details
// @Description Gets git-token provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderGitToken
// @Param   id path string true "provider ID"
// @Success 200 {object} gittokendto.GetGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/git-tokens/{id} [get]
func (h *ProvidersHandler) GetGitToken(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeGitToken, base.SettingScopeGlobal)
}

// CreateGitToken Creates a new git-token provider
// @Summary Creates a new git-token provider
// @Description Creates a new git-token provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderGitToken
// @Param   body body gittokendto.CreateGitTokenReq true "request data"
// @Success 201 {object} gittokendto.CreateGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/git-tokens [post]
func (h *ProvidersHandler) CreateGitToken(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeGitToken, base.SettingScopeGlobal)
}

// UpdateGitToken Updates git-token
// @Summary Updates git-token
// @Description Updates git-token
// @Tags    global_providers
// @Produce json
// @Id      updateProviderGitToken
// @Param   id path string true "provider ID"
// @Param   body body gittokendto.UpdateGitTokenReq true "request data"
// @Success 200 {object} gittokendto.UpdateGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/git-tokens/{id} [put]
func (h *ProvidersHandler) UpdateGitToken(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeGitToken, base.SettingScopeGlobal)
}

// UpdateGitTokenMeta Updates git-token meta
// @Summary Updates git-token meta
// @Description Updates git-token meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderGitTokenMeta
// @Param   id path string true "provider ID"
// @Param   body body gittokendto.UpdateGitTokenMetaReq true "request data"
// @Success 200 {object} gittokendto.UpdateGitTokenMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/git-tokens/{id}/meta [put]
func (h *ProvidersHandler) UpdateGitTokenMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeGitToken, base.SettingScopeGlobal)
}

// DeleteGitToken Deletes git-token provider
// @Summary Deletes git-token provider
// @Description Deletes git-token provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderGitToken
// @Param   id path string true "provider ID"
// @Success 200 {object} gittokendto.DeleteGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/git-tokens/{id} [delete]
func (h *ProvidersHandler) DeleteGitToken(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeGitToken, base.SettingScopeGlobal)
}

// TestGitTokenConn Test git-token connection
// @Summary Test git-token connection
// @Description Test git-token connection
// @Tags    global_providers
// @Produce json
// @Id      testGitTokenConn
// @Param   body body gittokendto.TestGitTokenConnReq true "request data"
// @Success 200 {object} gittokendto.TestGitTokenConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/git-tokens/test-conn [post]
func (h *ProvidersHandler) TestGitTokenConn(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := gittokendto.NewTestGitTokenConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GitTokenUC.TestGitTokenConn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
