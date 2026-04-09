package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuilduc/imagebuilddto"
)

// GetUniqueImageBuild Gets image build setting details
// @Summary Gets image build setting details
// @Description Gets image build setting details
// @Tags    project_settings
// @Produce json
// @Id      getProjectImageBuild
// @Param   projectID path string true "project ID"
// @Success 200 {object} imagebuilddto.GetUniqueImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build [get]
func (h *Handler) GetUniqueImageBuild(ctx *gin.Context) {
	h.GetUniqueSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}

// UpdateUniqueImageBuild Updates image build
// @Summary Updates image build
// @Description Updates image build
// @Tags    project_settings
// @Produce json
// @Id      updateProjectImageBuild
// @Param   projectID path string true "project ID"
// @Param   body body imagebuilddto.UpdateUniqueImageBuildReq true "request data"
// @Success 200 {object} imagebuilddto.UpdateUniqueImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build [put]
func (h *Handler) UpdateUniqueImageBuild(ctx *gin.Context) {
	h.UpdateUniqueSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}

// UpdateUniqueImageBuildMeta Updates image build meta
// @Summary Updates image build meta
// @Description Updates image build meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectImageBuildMeta
// @Param   projectID path string true "project ID"
// @Param   body body imagebuilddto.UpdateUniqueImageBuildMetaReq true "request data"
// @Success 200 {object} imagebuilddto.UpdateUniqueImageBuildMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build/meta [put]
func (h *Handler) UpdateUniqueImageBuildMeta(ctx *gin.Context) {
	h.UpdateUniqueSettingMeta(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}

// DeleteUniqueImageBuild Deletes image build setting
// @Summary Deletes image build setting
// @Description Deletes image build setting
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectImageBuild
// @Param   projectID path string true "project ID"
// @Success 200 {object} imagebuilddto.DeleteUniqueImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build [delete]
func (h *Handler) DeleteUniqueImageBuild(ctx *gin.Context) {
	h.DeleteUniqueSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}
