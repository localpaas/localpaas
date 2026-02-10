package providershandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
)

// ListBasicAuth Lists basic auth providers
// @Summary Lists basic auth providers
// @Description Lists basic auth providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderBasicAuth
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} basicauthdto.ListBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth [get]
func (h *ProvidersHandler) ListBasicAuth(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeGlobal)
}

// GetBasicAuth Gets basic auth provider details
// @Summary Gets basic auth provider details
// @Description Gets basic auth provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderBasicAuth
// @Param   itemID path string true "setting ID"
// @Success 200 {object} basicauthdto.GetBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth/{itemID} [get]
func (h *ProvidersHandler) GetBasicAuth(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeGlobal)
}

// CreateBasicAuth Creates a new basic auth provider
// @Summary Creates a new basic auth provider
// @Description Creates a new basic auth provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderBasicAuth
// @Param   body body basicauthdto.CreateBasicAuthReq true "request data"
// @Success 201 {object} basicauthdto.CreateBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth [post]
func (h *ProvidersHandler) CreateBasicAuth(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeGlobal)
}

// UpdateBasicAuth Updates basic auth
// @Summary Updates basic auth
// @Description Updates basic auth
// @Tags    global_providers
// @Produce json
// @Id      updateProviderBasicAuth
// @Param   itemID path string true "setting ID"
// @Param   body body basicauthdto.UpdateBasicAuthReq true "request data"
// @Success 200 {object} basicauthdto.UpdateBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth/{itemID} [put]
func (h *ProvidersHandler) UpdateBasicAuth(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeGlobal)
}

// UpdateBasicAuthMeta Updates basic auth meta
// @Summary Updates basic auth meta
// @Description Updates basic auth meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderBasicAuthMeta
// @Param   itemID path string true "setting ID"
// @Param   body body basicauthdto.UpdateBasicAuthMetaReq true "request data"
// @Success 200 {object} basicauthdto.UpdateBasicAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth/{itemID}/meta [put]
func (h *ProvidersHandler) UpdateBasicAuthMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeBasicAuth, base.SettingScopeGlobal)
}

// DeleteBasicAuth Deletes basic auth provider
// @Summary Deletes basic auth provider
// @Description Deletes basic auth provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderBasicAuth
// @Param   itemID path string true "setting ID"
// @Success 200 {object} basicauthdto.DeleteBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth/{itemID} [delete]
func (h *ProvidersHandler) DeleteBasicAuth(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeGlobal)
}
