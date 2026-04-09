package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuilduc/imagebuilddto"
)

// GetUniqueImageBuild Gets image build setting details
// @Summary Gets image build setting details
// @Description Gets image build setting details
// @Tags    settings
// @Produce json
// @Id      getSettingImageBuild
// @Success 200 {object} imagebuilddto.GetUniqueImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build [get]
func (h *Handler) GetUniqueImageBuild(ctx *gin.Context) {
	h.GetUniqueSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}

// UpdateUniqueImageBuild Updates image build
// @Summary Updates image build
// @Description Updates image build
// @Tags    settings
// @Produce json
// @Id      updateSettingImageBuild
// @Param   body body imagebuilddto.UpdateUniqueImageBuildReq true "request data"
// @Success 200 {object} imagebuilddto.UpdateUniqueImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build [put]
func (h *Handler) UpdateUniqueImageBuild(ctx *gin.Context) {
	h.UpdateUniqueSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}

// UpdateUniqueImageBuildMeta Updates image build meta
// @Summary Updates image build meta
// @Description Updates image build meta
// @Tags    settings
// @Produce json
// @Id      updateSettingImageBuildMeta
// @Param   body body imagebuilddto.UpdateUniqueImageBuildMetaReq true "request data"
// @Success 200 {object} imagebuilddto.UpdateUniqueImageBuildMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build/meta [put]
func (h *Handler) UpdateUniqueImageBuildMeta(ctx *gin.Context) {
	h.UpdateUniqueSettingMeta(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}

// DeleteUniqueImageBuild Deletes image build setting
// @Summary Deletes image build setting
// @Description Deletes image build setting
// @Tags    settings
// @Produce json
// @Id      deleteSettingImageBuild
// @Success 200 {object} imagebuilddto.DeleteUniqueImageBuildResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/image-build [delete]
func (h *Handler) DeleteUniqueImageBuild(ctx *gin.Context) {
	h.DeleteUniqueSetting(ctx, base.ResourceTypeImageBuild, base.SettingScopeGlobal)
}
