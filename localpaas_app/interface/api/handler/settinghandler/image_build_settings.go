package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuildsettingsuc/imagebuildsettingsdto"
)

// GetUniqueImageBuildSettings Gets image build settings
// @Summary Gets image build settings
// @Description Gets image build settings
// @Tags    settings
// @Produce json
// @Id      getSettingImageBuildSettings
// @Success 200 {object} imagebuildsettingsdto.GetUniqueImageBuildSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build-settings [get]
func (h *Handler) GetUniqueImageBuildSettings(ctx *gin.Context) {
	h.GetUniqueSetting(ctx, base.ResourceTypeImageBuildSettings, base.SettingScopeGlobal)
}

// UpdateUniqueImageBuildSettings Updates image build settings
// @Summary Updates image build settings
// @Description Updates image build settings
// @Tags    settings
// @Produce json
// @Id      updateSettingImageBuildSettings
// @Param   body body imagebuildsettingsdto.UpdateUniqueImageBuildSettingsReq true "request data"
// @Success 200 {object} imagebuildsettingsdto.UpdateUniqueImageBuildSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build-settings [put]
func (h *Handler) UpdateUniqueImageBuildSettings(ctx *gin.Context) {
	h.UpdateUniqueSetting(ctx, base.ResourceTypeImageBuildSettings, base.SettingScopeGlobal)
}

// UpdateUniqueImageBuildSettingsMeta Updates image build settings meta
// @Summary Updates image build settings meta
// @Description Updates image build settings meta
// @Tags    settings
// @Produce json
// @Id      updateSettingImageBuildSettingsMeta
// @Param   body body imagebuildsettingsdto.UpdateUniqueImageBuildSettingsMetaReq true "request data"
// @Success 200 {object} imagebuildsettingsdto.UpdateUniqueImageBuildSettingsMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build-settings/meta [put]
func (h *Handler) UpdateUniqueImageBuildSettingsMeta(ctx *gin.Context) {
	h.UpdateUniqueSettingMeta(ctx, base.ResourceTypeImageBuildSettings, base.SettingScopeGlobal)
}

// DeleteUniqueImageBuildSettings Deletes image build settings
// @Summary Deletes image build settings
// @Description Deletes image build settings
// @Tags    settings
// @Produce json
// @Id      deleteSettingImageBuildSettings
// @Success 200 {object} imagebuildsettingsdto.DeleteUniqueImageBuildSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build-settings [delete]
func (h *Handler) DeleteUniqueImageBuildSettings(ctx *gin.Context) {
	h.DeleteUniqueSetting(ctx, base.ResourceTypeImageBuildSettings, base.SettingScopeGlobal)
}
