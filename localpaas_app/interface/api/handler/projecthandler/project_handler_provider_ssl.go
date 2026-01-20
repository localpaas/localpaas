package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListSsl Lists SSL providers
// @Summary Lists SSL providers
// @Description Lists SSL providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectSSL
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} ssldto.ListSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls [get]
func (h *ProjectHandler) ListSsl(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewListSslReq()
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sslUC.ListSsl(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetSsl Gets SSL provider details
// @Summary Gets SSL provider details
// @Description Gets SSL provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectSSL
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} ssldto.GetSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls/{id} [get]
func (h *ProjectHandler) GetSsl(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewGetSslReq()
	req.ID = id
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sslUC.GetSsl(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateSsl Creates a new SSL provider
// @Summary Creates a new SSL provider
// @Description Creates a new SSL provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectSSL
// @Param   projectID path string true "project ID"
// @Param   body body ssldto.CreateSslReq true "request data"
// @Success 201 {object} ssldto.CreateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls [post]
func (h *ProjectHandler) CreateSsl(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewCreateSslReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sslUC.CreateSsl(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateSsl Updates SSL
// @Summary Updates SSL
// @Description Updates SSL
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSSL
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body ssldto.UpdateSslReq true "request data"
// @Success 200 {object} ssldto.UpdateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls/{id} [put]
func (h *ProjectHandler) UpdateSsl(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewUpdateSslReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sslUC.UpdateSsl(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateSslMeta Updates SSL meta
// @Summary Updates SSL meta
// @Description Updates SSL meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSSLMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body ssldto.UpdateSslMetaReq true "request data"
// @Success 200 {object} ssldto.UpdateSslMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls/{id}/meta [put]
func (h *ProjectHandler) UpdateSslMeta(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewUpdateSslMetaReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sslUC.UpdateSslMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteSsl Deletes SSL provider
// @Summary Deletes SSL provider
// @Description Deletes SSL provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectSSL
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} ssldto.DeleteSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls/{id} [delete]
func (h *ProjectHandler) DeleteSsl(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewDeleteSslReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sslUC.DeleteSsl(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
