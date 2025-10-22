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

// GetProjectEnvEnvVars Gets project env's env vars
// @Summary Gets project env's env vars
// @Description Gets project env's env vars
// @Settings    project_envs_env_vars
// @Produce json
// @Id      getProjectEnvEnvVars
// @Param   projectID path string true "project ID"
// @Param   projectEnvID path string true "project env ID"
// @Param   body body projectenvdto.GetProjectEnvEnvVarsReq true "request data"
// @Success 200 {object} projectenvdto.GetProjectEnvEnvVarsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/envs/{projectEnvID}/env-vars [get]
func (h *ProjectEnvHandler) GetProjectEnvEnvVars(ctx *gin.Context) {
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

	req := projectenvdto.NewGetProjectEnvEnvVarsReq()
	req.ProjectID = projectID
	req.ProjectEnvID = projectEnvID
	if err := h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectEnvUC.GetProjectEnvEnvVars(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateProjectEnvEnvVars Updates project env's env vars
// @Summary Updates project env's env vars
// @Description Updates project env's env vars
// @Settings    project_envs_env_vars
// @Produce json
// @Id      updateProjectEnvEnvVars
// @Param   projectID path string true "project ID"
// @Param   body body projectenvdto.UpdateProjectEnvEnvVarsReq true "request data"
// @Success 200 {object} projectenvdto.UpdateProjectEnvEnvVarsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/envs/{projectEnvID}/env-vars [put]
func (h *ProjectEnvHandler) UpdateProjectEnvEnvVars(ctx *gin.Context) {
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

	req := projectenvdto.NewUpdateProjectEnvEnvVarsReq()
	req.ProjectID = projectID
	req.ProjectEnvID = projectEnvID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectEnvUC.UpdateProjectEnvEnvVars(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
