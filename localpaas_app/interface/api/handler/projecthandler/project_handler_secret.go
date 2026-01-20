package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/secretuc/secretdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListProjectSecrets Lists project secrets
// @Summary Lists project secrets
// @Description Lists project secrets
// @Tags    projects
// @Produce json
// @Id      listProjectSecret
// @Param   projectID path string true "project ID"
// @Param   type query string false "`type=<setting type>`"
// @Success 200 {object} secretdto.ListSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/secrets [get]
func (h *ProjectHandler) ListProjectSecrets(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewListSecretReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.secretUC.ListSecret(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateProjectSecret Creates a project secret
// @Summary Creates a project secret
// @Description Creates a project secret
// @Tags    projects
// @Produce json
// @Id      createProjectSecret
// @Param   projectID path string true "project ID"
// @Param   body body secretdto.CreateSecretReq true "request data"
// @Success 201 {object} secretdto.CreateSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/secrets [post]
func (h *ProjectHandler) CreateProjectSecret(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewCreateSecretReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.secretUC.CreateSecret(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteProjectSecret Deletes a project secret
// @Summary Deletes a project secret
// @Description Deletes a project secret
// @Tags    projects
// @Produce json
// @Id      deleteProjectSecret
// @Param   projectID path string true "project ID"
// @Param   id path string true "secret ID"
// @Success 200 {object} secretdto.DeleteSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/secrets/{id} [delete]
func (h *ProjectHandler) DeleteProjectSecret(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewDeleteSecretReq()
	req.ID = id
	req.ProjectID = projectID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.secretUC.DeleteSecret(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
