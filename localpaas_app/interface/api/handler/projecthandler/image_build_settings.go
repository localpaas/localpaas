package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuildsettingsuc/imagebuildsettingsdto"
)

// GetUniqueImageBuildSettings Gets image build setting details
// @Summary Gets image build setting details
// @Description Gets image build setting details
// @Tags    project_settings
// @Produce json
// @Id      getProjectImageBuildSettings
// @Param   projectID path string true "project ID"
// @Success 200 {object} imagebuildsettingsdto.GetUniqueImageBuildSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build-settings [get]
func (h *Handler) GetUniqueImageBuildSettings(ctx *gin.Context) {
	h.GetUniqueSetting(ctx, base.ResourceTypeImageBuildSettings, base.SettingScopeProject)
}

// UpdateUniqueImageBuildSettings Updates image build settings
// @Summary Updates image build settings
// @Description Updates image build settings
// @Tags    project_settings
// @Produce json
// @Id      updateProjectImageBuildSettings
// @Param   projectID path string true "project ID"
// @Param   body body imagebuildsettingsdto.UpdateUniqueImageBuildSettingsReq true "request data"
// @Success 200 {object} imagebuildsettingsdto.UpdateUniqueImageBuildSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build-settings [put]
func (h *Handler) UpdateUniqueImageBuildSettings(ctx *gin.Context) {
	h.UpdateUniqueSetting(ctx, base.ResourceTypeImageBuildSettings, base.SettingScopeProject)
}

// UpdateUniqueImageBuildSettingsMeta Updates image build meta
// @Summary Updates image build meta
// @Description Updates image build meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectImageBuildSettingsMeta
// @Param   projectID path string true "project ID"
// @Param   body body imagebuildsettingsdto.UpdateUniqueImageBuildSettingsMetaReq true "request data"
// @Success 200 {object} imagebuildsettingsdto.UpdateUniqueImageBuildSettingsMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build-settings/meta [put]
func (h *Handler) UpdateUniqueImageBuildSettingsMeta(ctx *gin.Context) {
	h.UpdateUniqueSettingMeta(ctx, base.ResourceTypeImageBuildSettings, base.SettingScopeProject)
}

// DeleteUniqueImageBuildSettings Deletes image build settings
// @Summary Deletes image build settings
// @Description Deletes image build settings
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectImageBuildSettings
// @Param   projectID path string true "project ID"
// @Success 200 {object} imagebuildsettingsdto.DeleteUniqueImageBuildSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build-settings [delete]
func (h *Handler) DeleteUniqueImageBuildSettings(ctx *gin.Context) {
	h.DeleteUniqueSetting(ctx, base.ResourceTypeImageBuildSettings, base.SettingScopeProject)
}
