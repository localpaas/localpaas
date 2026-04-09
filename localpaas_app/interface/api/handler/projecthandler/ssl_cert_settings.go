package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/sslcertsettingsuc/sslcertsettingsdto"
)

// GetUniqueSSLCertSettings Gets ssl cert settings details
// @Summary Gets ssl cert settings details
// @Description Gets ssl cert settings details
// @Tags    project_settings
// @Produce json
// @Id      getProjectSSLCertSettings
// @Param   projectID path string true "project ID"
// @Success 200 {object} sslcertsettingsdto.GetUniqueSSLCertSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-cert-settings [get]
func (h *Handler) GetUniqueSSLCertSettings(ctx *gin.Context) {
	h.GetUniqueSetting(ctx, base.ResourceTypeSSLCertSettings, base.SettingScopeProject)
}

// UpdateUniqueSSLCertSettings Updates ssl cert settings
// @Summary Updates ssl cert settings
// @Description Updates ssl cert settings
// @Tags    project_settings
// @Produce json
// @Id      updateProjectSSLCertSettings
// @Param   projectID path string true "project ID"
// @Param   body body sslcertsettingsdto.UpdateUniqueSSLCertSettingsReq true "request data"
// @Success 200 {object} sslcertsettingsdto.UpdateUniqueSSLCertSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-cert-settings [put]
func (h *Handler) UpdateUniqueSSLCertSettings(ctx *gin.Context) {
	h.UpdateUniqueSetting(ctx, base.ResourceTypeSSLCertSettings, base.SettingScopeProject)
}

// UpdateUniqueSSLCertSettingsMeta Updates ssl cert settings meta
// @Summary Updates ssl cert settings meta
// @Description Updates ssl cert settings meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectSSLCertSettingsMeta
// @Param   projectID path string true "project ID"
// @Param   body body sslcertsettingsdto.UpdateUniqueSSLCertSettingsMetaReq true "request data"
// @Success 200 {object} sslcertsettingsdto.UpdateUniqueSSLCertSettingsMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-cert-settings/meta [put]
func (h *Handler) UpdateUniqueSSLCertSettingsMeta(ctx *gin.Context) {
	h.UpdateUniqueSettingMeta(ctx, base.ResourceTypeSSLCertSettings, base.SettingScopeProject)
}

// DeleteUniqueSSLCertSettings Deletes ssl cert settings setting
// @Summary Deletes ssl cert settings setting
// @Description Deletes ssl cert settings setting
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectSSLCertSettings
// @Param   projectID path string true "project ID"
// @Success 200 {object} sslcertsettingsdto.DeleteUniqueSSLCertSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-cert-settings [delete]
func (h *Handler) DeleteUniqueSSLCertSettings(ctx *gin.Context) {
	h.DeleteUniqueSetting(ctx, base.ResourceTypeSSLCertSettings, base.SettingScopeProject)
}
