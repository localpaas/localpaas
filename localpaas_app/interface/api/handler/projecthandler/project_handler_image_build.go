package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuilduc/imagebuilddto"
)

// ListImageBuild Lists image build settings
// @Summary Lists image build settings
// @Description Lists image build settings
// @Tags    project_settings
// @Produce json
// @Id      listProjectImageBuild
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} imagebuilddto.ListImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build [get]
func (h *ProjectHandler) ListImageBuild(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}

// GetImageBuild Gets image build setting details
// @Summary Gets image build setting details
// @Description Gets image build setting details
// @Tags    project_settings
// @Produce json
// @Id      getProjectImageBuild
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} imagebuilddto.GetImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build/{itemID} [get]
func (h *ProjectHandler) GetImageBuild(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}

// CreateImageBuild Creates a new image build setting
// @Summary Creates a new image build setting
// @Description Creates a new image build setting
// @Tags    project_settings
// @Produce json
// @Id      createProjectImageBuild
// @Param   projectID path string true "project ID"
// @Param   body body imagebuilddto.CreateImageBuildReq true "request data"
// @Success 201 {object} imagebuilddto.CreateImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build [post]
func (h *ProjectHandler) CreateImageBuild(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}

// UpdateImageBuild Updates image build
// @Summary Updates image build
// @Description Updates image build
// @Tags    project_settings
// @Produce json
// @Id      updateProjectImageBuild
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body imagebuilddto.UpdateImageBuildReq true "request data"
// @Success 200 {object} imagebuilddto.UpdateImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build/{itemID} [put]
func (h *ProjectHandler) UpdateImageBuild(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}

// UpdateImageBuildMeta Updates image build meta
// @Summary Updates image build meta
// @Description Updates image build meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectImageBuildMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body imagebuilddto.UpdateImageBuildMetaReq true "request data"
// @Success 200 {object} imagebuilddto.UpdateImageBuildMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build/{itemID}/meta [put]
func (h *ProjectHandler) UpdateImageBuildMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}

// DeleteImageBuild Deletes image build setting
// @Summary Deletes image build setting
// @Description Deletes image build setting
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectImageBuild
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} imagebuilddto.DeleteImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/image-build/{itemID} [delete]
func (h *ProjectHandler) DeleteImageBuild(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeProject)
}
