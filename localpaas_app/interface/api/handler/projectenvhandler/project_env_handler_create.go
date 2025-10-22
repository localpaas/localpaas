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

// CreateProjectEnv Creates a new project env
// @Summary Creates a new project env
// @Description Creates a new project env
// @Env     projects_envs
// @Produce json
// @Id      createProjectEnv
// @Param   projectID path string true "project ID"
// @Param   body body projectenvdto.CreateProjectEnvReq true "request data"
// @Success 201 {object} projectenvdto.CreateProjectEnvResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/envs [post]
func (h *ProjectEnvHandler) CreateProjectEnv(ctx *gin.Context) {
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

	req := projectenvdto.NewCreateProjectEnvReq()
	req.ProjectID = projectID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectEnvUC.CreateProjectEnv(h.RequestCtx(ctx), auth, req)
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
// @Param   body body projectenvdto.DeleteProjectEnvReq true "request data"
// @Success 200 {object} projectenvdto.DeleteProjectEnvResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/envs/{projectEnvID} [delete]
func (h *ProjectEnvHandler) DeleteProjectEnv(ctx *gin.Context) {
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

	req := projectenvdto.NewDeleteProjectEnvReq()
	req.ProjectID = projectID
	req.ProjectEnvID = projectEnvID
	if err := h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectEnvUC.DeleteProjectEnv(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
