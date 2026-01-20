package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListBasicAuth Lists basic auth providers
// @Summary Lists basic auth providers
// @Description Lists basic auth providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderBasicAuth
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} basicauthdto.ListBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth [get]
func (h *ProvidersHandler) ListBasicAuth(ctx *gin.Context) {
	auth, _, err := h.getAuth(ctx, base.ResourceTypeBasicAuth, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewListBasicAuthReq()
	req.GlobalOnly = true
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
// @Tags    global_providers
// @Produce json
// @Id      getProviderBasicAuth
// @Param   id path string true "provider ID"
// @Success 200 {object} basicauthdto.GetBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth/{id} [get]
func (h *ProvidersHandler) GetBasicAuth(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeBasicAuth, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewGetBasicAuthReq()
	req.ID = id
	req.GlobalOnly = true
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
// @Tags    global_providers
// @Produce json
// @Id      createProviderBasicAuth
// @Param   body body basicauthdto.CreateBasicAuthReq true "request data"
// @Success 201 {object} basicauthdto.CreateBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth [post]
func (h *ProvidersHandler) CreateBasicAuth(ctx *gin.Context) {
	auth, _, err := h.getAuth(ctx, base.ResourceTypeBasicAuth, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewCreateBasicAuthReq()
	req.GlobalOnly = true
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
// @Tags    global_providers
// @Produce json
// @Id      updateProviderBasicAuth
// @Param   id path string true "provider ID"
// @Param   body body basicauthdto.UpdateBasicAuthReq true "request data"
// @Success 200 {object} basicauthdto.UpdateBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth/{id} [put]
func (h *ProvidersHandler) UpdateBasicAuth(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeBasicAuth, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewUpdateBasicAuthReq()
	req.ID = id
	req.GlobalOnly = true
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
// @Tags    global_providers
// @Produce json
// @Id      updateProviderBasicAuthMeta
// @Param   id path string true "provider ID"
// @Param   body body basicauthdto.UpdateBasicAuthMetaReq true "request data"
// @Success 200 {object} basicauthdto.UpdateBasicAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth/{id}/meta [put]
func (h *ProvidersHandler) UpdateBasicAuthMeta(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeBasicAuth, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewUpdateBasicAuthMetaReq()
	req.ID = id
	req.GlobalOnly = true
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
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderBasicAuth
// @Param   id path string true "provider ID"
// @Success 200 {object} basicauthdto.DeleteBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/basic-auth/{id} [delete]
func (h *ProvidersHandler) DeleteBasicAuth(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeBasicAuth, base.ActionTypeDelete, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := basicauthdto.NewDeleteBasicAuthReq()
	req.ID = id
	req.GlobalOnly = true
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
