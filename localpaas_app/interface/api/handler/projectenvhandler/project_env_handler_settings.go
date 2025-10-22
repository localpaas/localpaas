package projectenvhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// GetProjectEnvSettings Gets project settings
// @Summary Gets project settings
// @Description Gets project settings
// @Settings    project_envs_settings
// @Produce json
// @Id      getProjectEnvSettings
// @Param   projectID path string true "project ID"
// @Param   body body projectenvdto.GetProjectEnvSettingsReq true "request data"
// @Success 200 {object} projectenvdto.GetProjectEnvSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/envs/{projectEnvID}/settings [get]
func (h *ProjectEnvHandler) GetProjectEnvSettings(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	projectEnvID, err := h.ParseStringParam(ctx, "projectEnvID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeProject,
		ResourceID:   projectID,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := projectenvdto.NewGetProjectEnvSettingsReq()
	req.ProjectID = projectID
	req.ProjectEnvID = projectEnvID
	if err := h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectEnvUC.GetProjectEnvSettings(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateProjectEnvSettings Updates project env settings
// @Summary Updates project env settings
// @Description Updates project env settings
// @Settings    project_envs_settings
// @Produce json
// @Id      updateProjectEnvSettings
// @Param   projectID path string true "project ID"
// @Param   projectEnvID path string true "project env ID"
// @Param   body body projectenvdto.UpdateProjectEnvSettingsReq true "request data"
// @Success 200 {object} projectenvdto.UpdateProjectEnvSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/envs/{projectEnvID}/settings [put]
func (h *ProjectEnvHandler) UpdateProjectEnvSettings(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	projectEnvID, err := h.ParseStringParam(ctx, "projectEnvID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeProject,
		ResourceID:   projectID,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := projectenvdto.NewUpdateProjectEnvSettingsReq()
	req.ProjectID = projectID
	req.ProjectEnvID = projectEnvID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectEnvUC.UpdateProjectEnvSettings(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
