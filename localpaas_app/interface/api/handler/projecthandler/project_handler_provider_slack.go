package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/slackuc/slackdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListSlack Lists Slack providers
// @Summary Lists Slack providers
// @Description Lists Slack providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectSlacks
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} slackdto.ListSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack [get]
func (h *ProjectHandler) ListSlack(ctx *gin.Context) {
	auth, projectID, _, err := h.getProjectProviderAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewListSlackReq()
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.ListSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetSlack Gets Slack provider details
// @Summary Gets Slack provider details
// @Description Gets Slack provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectSlack
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} slackdto.GetSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack/{id} [get]
func (h *ProjectHandler) GetSlack(ctx *gin.Context) {
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewGetSlackReq()
	req.ID = id
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.GetSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateSlack Creates a new Slack provider
// @Summary Creates a new Slack provider
// @Description Creates a new Slack provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectSlack
// @Param   projectID path string true "project ID"
// @Param   body body slackdto.CreateSlackReq true "request data"
// @Success 201 {object} slackdto.CreateSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack [post]
func (h *ProjectHandler) CreateSlack(ctx *gin.Context) {
	auth, projectID, _, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewCreateSlackReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.CreateSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateSlack Updates Slack provider
// @Summary Updates Slack provider
// @Description Updates Slack provider
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSlack
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body slackdto.UpdateSlackReq true "request data"
// @Success 200 {object} slackdto.UpdateSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack/{id} [put]
func (h *ProjectHandler) UpdateSlack(ctx *gin.Context) {
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewUpdateSlackReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.UpdateSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateSlackMeta Updates Slack meta provider
// @Summary Updates Slack meta provider
// @Description Updates Slack meta provider
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSlackMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body slackdto.UpdateSlackMetaReq true "request data"
// @Success 200 {object} slackdto.UpdateSlackMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack/{id}/meta [put]
func (h *ProjectHandler) UpdateSlackMeta(ctx *gin.Context) {
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewUpdateSlackMetaReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.UpdateSlackMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteSlack Deletes Slack provider
// @Summary Deletes Slack provider
// @Description Deletes Slack provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectSlack
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} slackdto.DeleteSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack/{id} [delete]
func (h *ProjectHandler) DeleteSlack(ctx *gin.Context) {
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewDeleteSlackReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.DeleteSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
