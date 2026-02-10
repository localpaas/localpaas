package providershandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
)

// ListOAuth Lists oauth providers
// @Summary Lists oauth providers
// @Description Lists oauth providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderOAuth
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} oauthdto.ListOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth [get]
func (h *ProvidersHandler) ListOAuth(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeOAuth, base.SettingScopeGlobal)
}

// GetOAuth Gets oauth provider details
// @Summary Gets oauth provider details
// @Description Gets oauth provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderOAuth
// @Param   itemID path string true "setting ID"
// @Success 200 {object} oauthdto.GetOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth/{itemID} [get]
func (h *ProvidersHandler) GetOAuth(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeOAuth, base.SettingScopeGlobal)
}

// CreateOAuth Creates a new oauth provider
// @Summary Creates a new oauth provider
// @Description Creates a new oauth provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderOAuth
// @Param   body body oauthdto.CreateOAuthReq true "request data"
// @Success 201 {object} oauthdto.CreateOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth [post]
func (h *ProvidersHandler) CreateOAuth(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeOAuth, base.SettingScopeGlobal)
}

// UpdateOAuth Updates oauth
// @Summary Updates oauth
// @Description Updates oauth
// @Tags    global_providers
// @Produce json
// @Id      updateProviderOAuth
// @Param   itemID path string true "setting ID"
// @Param   body body oauthdto.UpdateOAuthReq true "request data"
// @Success 200 {object} oauthdto.UpdateOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth/{itemID} [put]
func (h *ProvidersHandler) UpdateOAuth(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeOAuth, base.SettingScopeGlobal)
}

// UpdateOAuthMeta Updates oauth meta
// @Summary Updates oauth meta
// @Description Updates oauth meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderOAuthMeta
// @Param   itemID path string true "setting ID"
// @Param   body body oauthdto.UpdateOAuthMetaReq true "request data"
// @Success 200 {object} oauthdto.UpdateOAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth/{itemID}/meta [put]
func (h *ProvidersHandler) UpdateOAuthMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeOAuth, base.SettingScopeGlobal)
}

// DeleteOAuth Deletes oauth provider
// @Summary Deletes oauth provider
// @Description Deletes oauth provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderOAuth
// @Param   itemID path string true "setting ID"
// @Success 200 {object} oauthdto.DeleteOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth/{itemID} [delete]
func (h *ProvidersHandler) DeleteOAuth(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeOAuth, base.SettingScopeGlobal)
}
