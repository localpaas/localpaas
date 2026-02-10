package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/accesstokenuc/accesstokendto"
)

// ListAccessToken Lists access-token providers
// @Summary Lists access-token providers
// @Description Lists access-token providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderAccessToken
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} accesstokendto.ListAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/access-tokens [get]
func (h *ProvidersHandler) ListAccessToken(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeGlobal)
}

// GetAccessToken Gets access-token provider details
// @Summary Gets access-token provider details
// @Description Gets access-token provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderAccessToken
// @Param   itemID path string true "setting ID"
// @Success 200 {object} accesstokendto.GetAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/access-tokens/{itemID} [get]
func (h *ProvidersHandler) GetAccessToken(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeGlobal)
}

// CreateAccessToken Creates a new access-token provider
// @Summary Creates a new access-token provider
// @Description Creates a new access-token provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderAccessToken
// @Param   body body accesstokendto.CreateAccessTokenReq true "request data"
// @Success 201 {object} accesstokendto.CreateAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/access-tokens [post]
func (h *ProvidersHandler) CreateAccessToken(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeGlobal)
}

// UpdateAccessToken Updates access-token
// @Summary Updates access-token
// @Description Updates access-token
// @Tags    global_providers
// @Produce json
// @Id      updateProviderAccessToken
// @Param   itemID path string true "setting ID"
// @Param   body body accesstokendto.UpdateAccessTokenReq true "request data"
// @Success 200 {object} accesstokendto.UpdateAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/access-tokens/{itemID} [put]
func (h *ProvidersHandler) UpdateAccessToken(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeGlobal)
}

// UpdateAccessTokenMeta Updates access-token meta
// @Summary Updates access-token meta
// @Description Updates access-token meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderAccessTokenMeta
// @Param   itemID path string true "setting ID"
// @Param   body body accesstokendto.UpdateAccessTokenMetaReq true "request data"
// @Success 200 {object} accesstokendto.UpdateAccessTokenMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/access-tokens/{itemID}/meta [put]
func (h *ProvidersHandler) UpdateAccessTokenMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeAccessToken, base.SettingScopeGlobal)
}

// DeleteAccessToken Deletes access-token provider
// @Summary Deletes access-token provider
// @Description Deletes access-token provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderAccessToken
// @Param   itemID path string true "setting ID"
// @Success 200 {object} accesstokendto.DeleteAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/access-tokens/{itemID} [delete]
func (h *ProvidersHandler) DeleteAccessToken(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeGlobal)
}

// TestAccessTokenConn Test access-token connection
// @Summary Test access-token connection
// @Description Test access-token connection
// @Tags    global_providers
// @Produce json
// @Id      testAccessTokenConn
// @Param   body body accesstokendto.TestAccessTokenConnReq true "request data"
// @Success 200 {object} accesstokendto.TestAccessTokenConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/access-tokens/test-conn [post]
func (h *ProvidersHandler) TestAccessTokenConn(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := accesstokendto.NewTestAccessTokenConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.AccessTokenUC.TestAccessTokenConn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
