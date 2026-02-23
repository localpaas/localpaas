package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuilduc/imagebuilddto"
)

// ListImageBuild Lists image build settings
// @Summary Lists image build settings
// @Description Lists image build settings
// @Tags    settings
// @Produce json
// @Id      listSettingImageBuild
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} imagebuilddto.ListImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build [get]
func (h *SettingHandler) ListImageBuild(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}

// GetImageBuild Gets image build setting details
// @Summary Gets image build setting details
// @Description Gets image build setting details
// @Tags    settings
// @Produce json
// @Id      getSettingImageBuild
// @Param   itemID path string true "setting ID"
// @Success 200 {object} imagebuilddto.GetImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build/{itemID} [get]
func (h *SettingHandler) GetImageBuild(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}

// CreateImageBuild Creates a new image build setting
// @Summary Creates a new image build setting
// @Description Creates a new image build setting
// @Tags    settings
// @Produce json
// @Id      createSettingImageBuild
// @Param   body body imagebuilddto.CreateImageBuildReq true "request data"
// @Success 201 {object} imagebuilddto.CreateImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build [post]
func (h *SettingHandler) CreateImageBuild(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}

// UpdateImageBuild Updates image build
// @Summary Updates image build
// @Description Updates image build
// @Tags    settings
// @Produce json
// @Id      updateSettingImageBuild
// @Param   itemID path string true "setting ID"
// @Param   body body imagebuilddto.UpdateImageBuildReq true "request data"
// @Success 200 {object} imagebuilddto.UpdateImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build/{itemID} [put]
func (h *SettingHandler) UpdateImageBuild(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}

// UpdateImageBuildMeta Updates image build meta
// @Summary Updates image build meta
// @Description Updates image build meta
// @Tags    settings
// @Produce json
// @Id      updateSettingImageBuildMeta
// @Param   itemID path string true "setting ID"
// @Param   body body imagebuilddto.UpdateImageBuildMetaReq true "request data"
// @Success 200 {object} imagebuilddto.UpdateImageBuildMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build/{itemID}/meta [put]
func (h *SettingHandler) UpdateImageBuildMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}

// DeleteImageBuild Deletes image build setting
// @Summary Deletes image build setting
// @Description Deletes image build setting
// @Tags    settings
// @Produce json
// @Id      deleteSettingImageBuild
// @Param   itemID path string true "setting ID"
// @Success 200 {object} imagebuilddto.DeleteImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build/{itemID} [delete]
func (h *SettingHandler) DeleteImageBuild(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}
