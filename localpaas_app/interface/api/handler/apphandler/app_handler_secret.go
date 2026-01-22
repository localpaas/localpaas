package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
)

// ListAppSecrets Lists app secrets
// @Summary Lists app secrets
// @Description Lists app secrets
// @Tags    apps
// @Produce json
// @Id      listAppSecret
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} secretdto.ListSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/secrets [get]
func (h *AppHandler) ListAppSecrets(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewListSecretReq()
	req.ParentObjectID = projectID
	req.ObjectID = appID
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

// CreateAppSecret Creates an app secret
// @Summary Creates an app secret
// @Description Creates an app secret
// @Tags    apps
// @Produce json
// @Id      createAppSecret
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body secretdto.CreateSecretReq true "request data"
// @Success 201 {object} secretdto.CreateSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/secrets [post]
func (h *AppHandler) CreateAppSecret(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewCreateSecretReq()
	req.ParentObjectID = projectID
	req.ObjectID = appID
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

// DeleteAppSecret Deletes an app secret
// @Summary Deletes an app secret
// @Description Deletes an app secret
// @Tags    apps
// @Produce json
// @Id      deleteAppSecret
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "secret ID"
// @Success 200 {object} secretdto.DeleteSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/secrets/{id} [delete]
func (h *AppHandler) DeleteAppSecret(ctx *gin.Context) {
	auth, projectID, appID, itemID, err := h.getAuthForItem(ctx, base.ActionTypeWrite, "id")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewDeleteSecretReq()
	req.ID = itemID
	req.ParentObjectID = projectID
	req.ObjectID = appID
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
