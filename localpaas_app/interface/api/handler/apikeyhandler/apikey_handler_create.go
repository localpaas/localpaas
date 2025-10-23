package apikeyhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apikeyuc/apikeydto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// CreateAPIKey Creates a new API key
// @Summary Creates a new API key
// @Description Creates a new API key
// @Tags    api_keys
// @Produce json
// @Id      createAPIKey
// @Param   body body apikeydto.CreateAPIKeyReq true "request data"
// @Success 201 {object} apikeydto.CreateAPIKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /api-keys [post]
func (h *APIKeyHandler) CreateAPIKey(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeAPIKey,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := apikeydto.NewCreateAPIKeyReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.apiKeyUC.CreateAPIKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteAPIKey Deletes an API key
// @Summary Deletes an API key
// @Description Deletes an API key
// @Tags    api_keys
// @Produce json
// @Id      deleteAPIKey
// @Param   ID path string true "API key ID"
// @Success 200 {object} apikeydto.DeleteAPIKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /api-keys/{ID} [delete]
func (h *APIKeyHandler) DeleteAPIKey(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeAPIKey,
		ResourceID:   id,
		Action:       base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := apikeydto.NewDeleteAPIKeyReq()
	req.ID = id
	if err := h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.apiKeyUC.DeleteAPIKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
