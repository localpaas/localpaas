package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

// ListSSL Lists SSL settings
// @Summary Lists SSL settings
// @Description Lists SSL settings
// @Tags    settings
// @Produce json
// @Id      listSettingSSL
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} ssldto.ListSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls [get]
func (h *SettingHandler) ListSSL(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// GetSSL Gets SSL setting details
// @Summary Gets SSL setting details
// @Description Gets SSL setting details
// @Tags    settings
// @Produce json
// @Id      getSettingSSL
// @Param   itemID path string true "setting ID"
// @Success 200 {object} ssldto.GetSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls/{itemID} [get]
func (h *SettingHandler) GetSSL(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// CreateSSL Creates a new SSL setting
// @Summary Creates a new SSL setting
// @Description Creates a new SSL setting
// @Tags    settings
// @Produce json
// @Id      createSettingSSL
// @Param   body body ssldto.CreateSSLReq true "request data"
// @Success 201 {object} ssldto.CreateSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls [post]
func (h *SettingHandler) CreateSSL(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// UpdateSSL Updates SSL
// @Summary Updates SSL
// @Description Updates SSL
// @Tags    settings
// @Produce json
// @Id      updateSettingSSL
// @Param   itemID path string true "setting ID"
// @Param   body body ssldto.UpdateSSLReq true "request data"
// @Success 200 {object} ssldto.UpdateSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls/{itemID} [put]
func (h *SettingHandler) UpdateSSL(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// UpdateSSLMeta Updates SSL meta
// @Summary Updates SSL meta
// @Description Updates SSL meta
// @Tags    settings
// @Produce json
// @Id      updateSettingSSLMeta
// @Param   itemID path string true "setting ID"
// @Param   body body ssldto.UpdateSSLMetaReq true "request data"
// @Success 200 {object} ssldto.UpdateSSLMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls/{itemID}/meta [put]
func (h *SettingHandler) UpdateSSLMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}

// DeleteSSL Deletes SSL setting
// @Summary Deletes SSL setting
// @Description Deletes SSL setting
// @Tags    settings
// @Produce json
// @Id      deleteSettingSSL
// @Param   itemID path string true "setting ID"
// @Success 200 {object} ssldto.DeleteSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls/{itemID} [delete]
func (h *SettingHandler) DeleteSSL(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSSL, base.SettingScopeGlobal)
}
