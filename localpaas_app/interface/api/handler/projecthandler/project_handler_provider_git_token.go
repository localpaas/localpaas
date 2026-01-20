package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListGitToken Lists git-token providers
// @Summary Lists git-token providers
// @Description Lists git-token providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectGitTokens
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} gittokendto.ListGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/git-tokens [get]
func (h *ProjectHandler) ListGitToken(ctx *gin.Context) {
	auth, projectID, _, err := h.getProjectProviderAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := gittokendto.NewListGitTokenReq()
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.gitTokenUC.ListGitToken(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetGitToken Gets git-token provider details
// @Summary Gets git-token provider details
// @Description Gets git-token provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectGitToken
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} gittokendto.GetGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/git-tokens/{id} [get]
func (h *ProjectHandler) GetGitToken(ctx *gin.Context) {
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := gittokendto.NewGetGitTokenReq()
	req.ID = id
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.gitTokenUC.GetGitToken(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateGitToken Creates a new git-token provider
// @Summary Creates a new git-token provider
// @Description Creates a new git-token provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectGitToken
// @Param   projectID path string true "project ID"
// @Param   body body gittokendto.CreateGitTokenReq true "request data"
// @Success 201 {object} gittokendto.CreateGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/git-tokens [post]
func (h *ProjectHandler) CreateGitToken(ctx *gin.Context) {
	auth, projectID, _, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := gittokendto.NewCreateGitTokenReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.gitTokenUC.CreateGitToken(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateGitToken Updates git-token
// @Summary Updates git-token
// @Description Updates git-token
// @Tags    project_providers
// @Produce json
// @Id      updateProjectGitToken
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body gittokendto.UpdateGitTokenReq true "request data"
// @Success 200 {object} gittokendto.UpdateGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/git-tokens/{id} [put]
func (h *ProjectHandler) UpdateGitToken(ctx *gin.Context) {
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := gittokendto.NewUpdateGitTokenReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.gitTokenUC.UpdateGitToken(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateGitTokenMeta Updates git-token meta
// @Summary Updates git-token meta
// @Description Updates git-token meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectGitTokenMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body gittokendto.UpdateGitTokenMetaReq true "request data"
// @Success 200 {object} gittokendto.UpdateGitTokenMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/git-tokens/{id}/meta [put]
func (h *ProjectHandler) UpdateGitTokenMeta(ctx *gin.Context) {
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := gittokendto.NewUpdateGitTokenMetaReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.gitTokenUC.UpdateGitTokenMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteGitToken Deletes git-token provider
// @Summary Deletes git-token provider
// @Description Deletes git-token provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectGitToken
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} gittokendto.DeleteGitTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/git-tokens/{id} [delete]
func (h *ProjectHandler) DeleteGitToken(ctx *gin.Context) {
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := gittokendto.NewDeleteGitTokenReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.gitTokenUC.DeleteGitToken(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
