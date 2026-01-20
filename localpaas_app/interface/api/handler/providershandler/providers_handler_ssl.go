package providershandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

// ListSsl Lists SSL providers
// @Summary Lists SSL providers
// @Description Lists SSL providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderSSL
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} ssldto.ListSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls [get]
func (h *ProvidersHandler) ListSsl(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// GetSsl Gets SSL provider details
// @Summary Gets SSL provider details
// @Description Gets SSL provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderSSL
// @Param   id path string true "provider ID"
// @Success 200 {object} ssldto.GetSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{id} [get]
func (h *ProvidersHandler) GetSsl(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// CreateSsl Creates a new SSL provider
// @Summary Creates a new SSL provider
// @Description Creates a new SSL provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderSSL
// @Param   body body ssldto.CreateSslReq true "request data"
// @Success 201 {object} ssldto.CreateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls [post]
func (h *ProvidersHandler) CreateSsl(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// UpdateSsl Updates SSL
// @Summary Updates SSL
// @Description Updates SSL
// @Tags    global_providers
// @Produce json
// @Id      updateProviderSSL
// @Param   id path string true "provider ID"
// @Param   body body ssldto.UpdateSslReq true "request data"
// @Success 200 {object} ssldto.UpdateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{id} [put]
func (h *ProvidersHandler) UpdateSsl(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// UpdateSslMeta Updates SSL meta
// @Summary Updates SSL meta
// @Description Updates SSL meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderSSLMeta
// @Param   id path string true "provider ID"
// @Param   body body ssldto.UpdateSslMetaReq true "request data"
// @Success 200 {object} ssldto.UpdateSslMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{id}/meta [put]
func (h *ProvidersHandler) UpdateSslMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// DeleteSsl Deletes SSL provider
// @Summary Deletes SSL provider
// @Description Deletes SSL provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderSSL
// @Param   id path string true "provider ID"
// @Success 200 {object} ssldto.DeleteSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{id} [delete]
func (h *ProvidersHandler) DeleteSsl(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}
