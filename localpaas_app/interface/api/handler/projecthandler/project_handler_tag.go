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

// CreateProjectTag Creates a new project tag
// @Summary Creates a new project tag
// @Description Creates a new project tag
// @Tags    projects
// @Produce json
// @Id      createProjectTag
// @Param   projectID path string true "project ID"
// @Param   body body projectdto.CreateProjectTagReq true "request data"
// @Success 201 {object} projectdto.CreateProjectTagResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/tags [post]
func (h *ProjectHandler) CreateProjectTag(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProject,
		ResourceType:   base.ResourceTypeProject,
		ResourceID:     projectID,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := projectdto.NewCreateProjectTagReq()
	req.ProjectID = projectID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectUC.CreateProjectTag(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteProjectTags Deletes project tags
// @Summary Deletes project tags
// @Description Deletes project tags
// @Tags    projects
// @Produce json
// @Id      deleteProjectTags
// @Param   projectID path string true "project ID"
// @Param   body body projectdto.DeleteProjectTagsReq true "request data"
// @Success 200 {object} projectdto.DeleteProjectTagsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/tags/delete [post]
func (h *ProjectHandler) DeleteProjectTags(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProject,
		ResourceType:   base.ResourceTypeProject,
		ResourceID:     projectID,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := projectdto.NewDeleteProjectTagsReq()
	req.ProjectID = projectID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectUC.DeleteProjectTags(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
