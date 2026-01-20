package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListBasicAuth Lists basic auth providers
// @Summary Lists basic auth providers
// @Description Lists basic auth providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} basicauthdto.ListBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth [get]
func (h *ProjectHandler) ListBasicAuth(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewListBasicAuthReq()
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.basicAuthUC.ListBasicAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetBasicAuth Gets basic auth provider details
// @Summary Gets basic auth provider details
// @Description Gets basic auth provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} basicauthdto.GetBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth/{id} [get]
func (h *ProjectHandler) GetBasicAuth(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewGetBasicAuthReq()
	req.ID = id
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.basicAuthUC.GetBasicAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateBasicAuth Creates a new basic auth provider
// @Summary Creates a new basic auth provider
// @Description Creates a new basic auth provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   body body basicauthdto.CreateBasicAuthReq true "request data"
// @Success 201 {object} basicauthdto.CreateBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth [post]
func (h *ProjectHandler) CreateBasicAuth(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewCreateBasicAuthReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.basicAuthUC.CreateBasicAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateBasicAuth Updates basic auth
// @Summary Updates basic auth
// @Description Updates basic auth
// @Tags    project_providers
// @Produce json
// @Id      updateProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body basicauthdto.UpdateBasicAuthReq true "request data"
// @Success 200 {object} basicauthdto.UpdateBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth/{id} [put]
func (h *ProjectHandler) UpdateBasicAuth(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewUpdateBasicAuthReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.basicAuthUC.UpdateBasicAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateBasicAuthMeta Updates basic auth meta
// @Summary Updates basic auth meta
// @Description Updates basic auth meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectBasicAuthMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body basicauthdto.UpdateBasicAuthMetaReq true "request data"
// @Success 200 {object} basicauthdto.UpdateBasicAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth/{id}/meta [put]
func (h *ProjectHandler) UpdateBasicAuthMeta(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewUpdateBasicAuthMetaReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.basicAuthUC.UpdateBasicAuthMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteBasicAuth Deletes basic auth provider
// @Summary Deletes basic auth provider
// @Description Deletes basic auth provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} basicauthdto.DeleteBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth/{id} [delete]
func (h *ProjectHandler) DeleteBasicAuth(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewDeleteBasicAuthReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.basicAuthUC.DeleteBasicAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
