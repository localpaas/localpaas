package providershandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

// ListSSL Lists SSL providers
// @Summary Lists SSL providers
// @Description Lists SSL providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderSSL
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} ssldto.ListSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls [get]
func (h *ProvidersHandler) ListSSL(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// GetSSL Gets SSL provider details
// @Summary Gets SSL provider details
// @Description Gets SSL provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderSSL
// @Param   itemID path string true "setting ID"
// @Success 200 {object} ssldto.GetSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{itemID} [get]
func (h *ProvidersHandler) GetSSL(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// CreateSSL Creates a new SSL provider
// @Summary Creates a new SSL provider
// @Description Creates a new SSL provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderSSL
// @Param   body body ssldto.CreateSSLReq true "request data"
// @Success 201 {object} ssldto.CreateSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls [post]
func (h *ProvidersHandler) CreateSSL(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// UpdateSSL Updates SSL
// @Summary Updates SSL
// @Description Updates SSL
// @Tags    global_providers
// @Produce json
// @Id      updateProviderSSL
// @Param   itemID path string true "setting ID"
// @Param   body body ssldto.UpdateSSLReq true "request data"
// @Success 200 {object} ssldto.UpdateSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{itemID} [put]
func (h *ProvidersHandler) UpdateSSL(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// UpdateSSLMeta Updates SSL meta
// @Summary Updates SSL meta
// @Description Updates SSL meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderSSLMeta
// @Param   itemID path string true "setting ID"
// @Param   body body ssldto.UpdateSSLMetaReq true "request data"
// @Success 200 {object} ssldto.UpdateSSLMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{itemID}/meta [put]
func (h *ProvidersHandler) UpdateSSLMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// DeleteSSL Deletes SSL provider
// @Summary Deletes SSL provider
// @Description Deletes SSL provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderSSL
// @Param   itemID path string true "setting ID"
// @Success 200 {object} ssldto.DeleteSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{itemID} [delete]
func (h *ProvidersHandler) DeleteSSL(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}
