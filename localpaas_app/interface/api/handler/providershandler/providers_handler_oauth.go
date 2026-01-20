package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc/oauthdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListOAuth Lists oauth providers
// @Summary Lists oauth providers
// @Description Lists oauth providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderOAuth
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} oauthdto.ListOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth [get]
func (h *ProvidersHandler) ListOAuth(ctx *gin.Context) {
	auth, _, err := h.getAuth(ctx, base.ResourceTypeOAuth, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := oauthdto.NewListOAuthReq()
	req.GlobalOnly = true
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.oauthUC.ListOAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetOAuth Gets oauth provider details
// @Summary Gets oauth provider details
// @Description Gets oauth provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderOAuth
// @Param   id path string true "provider ID"
// @Success 200 {object} oauthdto.GetOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth/{id} [get]
func (h *ProvidersHandler) GetOAuth(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeOAuth, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := oauthdto.NewGetOAuthReq()
	req.ID = id
	req.GlobalOnly = true
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.oauthUC.GetOAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateOAuth Creates a new oauth provider
// @Summary Creates a new oauth provider
// @Description Creates a new oauth provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderOAuth
// @Param   body body oauthdto.CreateOAuthReq true "request data"
// @Success 201 {object} oauthdto.CreateOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth [post]
func (h *ProvidersHandler) CreateOAuth(ctx *gin.Context) {
	auth, _, err := h.getAuth(ctx, base.ResourceTypeOAuth, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := oauthdto.NewCreateOAuthReq()
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.oauthUC.CreateOAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateOAuth Updates oauth
// @Summary Updates oauth
// @Description Updates oauth
// @Tags    global_providers
// @Produce json
// @Id      updateProviderOAuth
// @Param   id path string true "provider ID"
// @Param   body body oauthdto.UpdateOAuthReq true "request data"
// @Success 200 {object} oauthdto.UpdateOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth/{id} [put]
func (h *ProvidersHandler) UpdateOAuth(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeOAuth, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := oauthdto.NewUpdateOAuthReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.oauthUC.UpdateOAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateOAuthMeta Updates oauth meta
// @Summary Updates oauth meta
// @Description Updates oauth meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderOAuthMeta
// @Param   id path string true "provider ID"
// @Param   body body oauthdto.UpdateOAuthMetaReq true "request data"
// @Success 200 {object} oauthdto.UpdateOAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth/{id}/meta [put]
func (h *ProvidersHandler) UpdateOAuthMeta(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeOAuth, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := oauthdto.NewUpdateOAuthMetaReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.oauthUC.UpdateOAuthMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteOAuth Deletes oauth provider
// @Summary Deletes oauth provider
// @Description Deletes oauth provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderOAuth
// @Param   id path string true "provider ID"
// @Success 200 {object} oauthdto.DeleteOAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/oauth/{id} [delete]
func (h *ProvidersHandler) DeleteOAuth(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeOAuth, base.ActionTypeDelete, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := oauthdto.NewDeleteOAuthReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.oauthUC.DeleteOAuth(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
