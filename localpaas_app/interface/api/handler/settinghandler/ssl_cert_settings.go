package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/sslcertsettingsuc/sslcertsettingsdto"
)

// GetUniqueSSLCertSettings Gets ssl cert settings details
// @Summary Gets ssl cert settings details
// @Description Gets ssl cert settings details
// @Tags    settings
// @Produce json
// @Id      getSettingSSLCertSettings
// @Success 200 {object} sslcertsettingsdto.GetUniqueSSLCertSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssl-cert-settings [get]
func (h *Handler) GetUniqueSSLCertSettings(ctx *gin.Context) {
	h.GetUniqueSetting(ctx, base.ResourceTypeSSLCertSettings, base.SettingScopeGlobal)
}

// UpdateUniqueSSLCertSettings Updates ssl cert settings
// @Summary Updates ssl cert settings
// @Description Updates ssl cert settings
// @Tags    settings
// @Produce json
// @Id      updateSettingSSLCertSettings
// @Param   body body sslcertsettingsdto.UpdateUniqueSSLCertSettingsReq true "request data"
// @Success 200 {object} sslcertsettingsdto.UpdateUniqueSSLCertSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssl-cert-settings [put]
func (h *Handler) UpdateUniqueSSLCertSettings(ctx *gin.Context) {
	h.UpdateUniqueSetting(ctx, base.ResourceTypeSSLCertSettings, base.SettingScopeGlobal)
}

// UpdateUniqueSSLCertSettingsMeta Updates ssl cert settings meta
// @Summary Updates ssl cert settings meta
// @Description Updates ssl cert settings meta
// @Tags    settings
// @Produce json
// @Id      updateSettingSSLCertSettingsMeta
// @Param   body body sslcertsettingsdto.UpdateUniqueSSLCertSettingsMetaReq true "request data"
// @Success 200 {object} sslcertsettingsdto.UpdateUniqueSSLCertSettingsMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssl-cert-settings/meta [put]
func (h *Handler) UpdateUniqueSSLCertSettingsMeta(ctx *gin.Context) {
	h.UpdateUniqueSettingMeta(ctx, base.ResourceTypeSSLCertSettings, base.SettingScopeGlobal)
}

// DeleteUniqueSSLCertSettings Deletes ssl cert settings setting
// @Summary Deletes ssl cert settings setting
// @Description Deletes ssl cert settings setting
// @Tags    settings
// @Produce json
// @Id      deleteSettingSSLCertSettings
// @Success 200 {object} sslcertsettingsdto.DeleteUniqueSSLCertSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssl-cert-settings [delete]
func (h *Handler) DeleteUniqueSSLCertSettings(ctx *gin.Context) {
	h.DeleteUniqueSetting(ctx, base.ResourceTypeSSLCertSettings, base.SettingScopeGlobal)
}
