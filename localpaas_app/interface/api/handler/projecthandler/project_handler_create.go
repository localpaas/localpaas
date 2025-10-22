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

// CreateProject Creates a new project
// @Summary Creates a new project
// @Description Creates a new project
// @Tags    projects
// @Produce json
// @Id      createProject
// @Param   body body projectdto.CreateProjectReq true "request data"
// @Success 201 {object} projectdto.CreateProjectResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects [post]
func (h *ProjectHandler) CreateProject(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeProject,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := projectdto.NewCreateProjectReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectUC.CreateProject(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteProject Deletes a project
// @Summary Deletes a project
// @Description Deletes a project
// @Tags    projects
// @Produce json
// @Id      deleteProject
// @Param   projectID path string true "project ID"
// @Success 200 {object} projectdto.DeleteProjectResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID} [delete]
func (h *ProjectHandler) DeleteProject(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeProject,
		ResourceID:   projectID,
		Action:       base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := projectdto.NewDeleteProjectReq()
	req.ProjectID = projectID
	if err := h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectUC.DeleteProject(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
