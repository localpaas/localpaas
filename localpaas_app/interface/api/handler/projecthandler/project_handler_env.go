package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// CreateProjectEnv Creates a new project tag
// @Summary Creates a new project tag
// @Description Creates a new project tag
// @Env     projects_envs
// @Produce json
// @Id      createProjectEnv
// @Param   projectID path string true "project ID"
// @Param   body body projectdto.CreateProjectEnvReq true "request data"
// @Success 201 {object} projectdto.CreateProjectEnvResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/envs [post]
func (h *ProjectHandler) CreateProjectEnv(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
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

	req := projectdto.NewCreateProjectEnvReq()
	req.ProjectID = projectID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectUC.CreateProjectEnv(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteProjectEnv Deletes a project env
// @Summary Deletes a project env
// @Description Deletes a project env
// @Env     projects_envs
// @Produce json
// @Id      deleteProjectEnv
// @Param   projectID path string true "project ID"
// @Param   projectEnvID path string true "project env ID"
// @Param   body body projectdto.DeleteProjectEnvReq true "request data"
// @Success 200 {object} projectdto.DeleteProjectEnvResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/envs/{projectEnvID} [delete]
func (h *ProjectHandler) DeleteProjectEnv(ctx *gin.Context) {
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

	req := projectdto.NewDeleteProjectEnvReq()
	req.ProjectID = projectID
	req.ProjectEnvID = projectEnvID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectUC.DeleteProjectEnv(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
