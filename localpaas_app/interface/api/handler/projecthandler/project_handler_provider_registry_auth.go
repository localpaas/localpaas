package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListRegistryAuth Lists registry auth providers
// @Summary Lists registry auth providers
// @Description Lists registry auth providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} registryauthdto.ListRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth [get]
func (h *ProjectHandler) ListRegistryAuth(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewListRegistryAuthReq()
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.registryAuthUC.ListRegistryAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetRegistryAuth Gets registry auth provider details
// @Summary Gets registry auth provider details
// @Description Gets registry auth provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} registryauthdto.GetRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth/{id} [get]
func (h *ProjectHandler) GetRegistryAuth(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewGetRegistryAuthReq()
	req.ID = id
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.registryAuthUC.GetRegistryAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateRegistryAuth Creates a new registry auth provider
// @Summary Creates a new registry auth provider
// @Description Creates a new registry auth provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   body body registryauthdto.CreateRegistryAuthReq true "request data"
// @Success 201 {object} registryauthdto.CreateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth [post]
func (h *ProjectHandler) CreateRegistryAuth(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewCreateRegistryAuthReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.registryAuthUC.CreateRegistryAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateRegistryAuth Updates registry auth
// @Summary Updates registry auth
// @Description Updates registry auth
// @Tags    project_providers
// @Produce json
// @Id      updateProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body registryauthdto.UpdateRegistryAuthReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth/{id} [put]
func (h *ProjectHandler) UpdateRegistryAuth(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewUpdateRegistryAuthReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.registryAuthUC.UpdateRegistryAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateRegistryAuthMeta Updates registry auth meta
// @Summary Updates registry auth meta
// @Description Updates registry auth meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectRegistryAuthMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body registryauthdto.UpdateRegistryAuthMetaReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth/{id}/meta [put]
func (h *ProjectHandler) UpdateRegistryAuthMeta(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewUpdateRegistryAuthMetaReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.registryAuthUC.UpdateRegistryAuthMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteRegistryAuth Deletes registry auth provider
// @Summary Deletes registry auth provider
// @Description Deletes registry auth provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} registryauthdto.DeleteRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth/{id} [delete]
func (h *ProjectHandler) DeleteRegistryAuth(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewDeleteRegistryAuthReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.registryAuthUC.DeleteRegistryAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
