package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

// ListRegistryAuth Lists registry auth providers
// @Summary Lists registry auth providers
// @Description Lists registry auth providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderRegistryAuth
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} registryauthdto.ListRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth [get]
func (h *ProvidersHandler) ListRegistryAuth(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeGlobal)
}

// GetRegistryAuth Gets registry auth provider details
// @Summary Gets registry auth provider details
// @Description Gets registry auth provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderRegistryAuth
// @Param   id path string true "provider ID"
// @Success 200 {object} registryauthdto.GetRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/{id} [get]
func (h *ProvidersHandler) GetRegistryAuth(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeGlobal)
}

// CreateRegistryAuth Creates a new registry auth provider
// @Summary Creates a new registry auth provider
// @Description Creates a new registry auth provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderRegistryAuth
// @Param   body body registryauthdto.CreateRegistryAuthReq true "request data"
// @Success 201 {object} registryauthdto.CreateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth [post]
func (h *ProvidersHandler) CreateRegistryAuth(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeGlobal)
}

// UpdateRegistryAuth Updates registry auth
// @Summary Updates registry auth
// @Description Updates registry auth
// @Tags    global_providers
// @Produce json
// @Id      updateProviderRegistryAuth
// @Param   id path string true "provider ID"
// @Param   body body registryauthdto.UpdateRegistryAuthReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/{id} [put]
func (h *ProvidersHandler) UpdateRegistryAuth(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeGlobal)
}

// UpdateRegistryAuthMeta Updates registry auth meta
// @Summary Updates registry auth meta
// @Description Updates registry auth meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderRegistryAuthMeta
// @Param   id path string true "provider ID"
// @Param   body body registryauthdto.UpdateRegistryAuthMetaReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/{id}/meta [put]
func (h *ProvidersHandler) UpdateRegistryAuthMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeGlobal)
}

// DeleteRegistryAuth Deletes registry auth provider
// @Summary Deletes registry auth provider
// @Description Deletes registry auth provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderRegistryAuth
// @Param   id path string true "provider ID"
// @Success 200 {object} registryauthdto.DeleteRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/{id} [delete]
func (h *ProvidersHandler) DeleteRegistryAuth(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeGlobal)
}

// TestRegistryAuthConn Tests registry auth connection
// @Summary Tests registry auth connection
// @Description Tests registry auth connection
// @Tags    global_providers
// @Produce json
// @Id      testRegistryAuthConn
// @Param   body body registryauthdto.TestRegistryAuthConnReq true "request data"
// @Success 200 {object} registryauthdto.TestRegistryAuthConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/test-conn [post]
func (h *ProvidersHandler) TestRegistryAuthConn(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewTestRegistryAuthConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.RegistryAuthUC.TestRegistryAuthConn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
