package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc/ssldto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListSsl Lists SSL providers
// @Summary Lists SSL providers
// @Description Lists SSL providers
// @Tags    providers_ssl
// @Produce json
// @Id      listSslProviders
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} ssldto.ListSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls [get]
func (h *ProvidersHandler) ListSsl(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSsl,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewListSslReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
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
// @Tags    providers_ssl
// @Produce json
// @Id      getSslProvider
// @Param   ID path string true "provider ID"
// @Success 200 {object} ssldto.GetSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{ID} [get]
func (h *ProvidersHandler) GetSsl(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSsl,
		ResourceID:     id,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewGetSslReq()
	req.ID = id
	if err = h.ParseRequest(ctx, req, nil); err != nil {
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
// @Tags    providers_ssl
// @Produce json
// @Id      createSslProvider
// @Param   body body ssldto.CreateSslReq true "request data"
// @Success 201 {object} ssldto.CreateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls [post]
func (h *ProvidersHandler) CreateSsl(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewCreateSslReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
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
// @Tags    providers_ssl
// @Produce json
// @Id      updateSslProvider
// @Param   ID path string true "provider ID"
// @Param   body body ssldto.UpdateSslReq true "request data"
// @Success 200 {object} ssldto.UpdateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{ID} [put]
func (h *ProvidersHandler) UpdateSsl(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSsl,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewUpdateSslReq()
	req.ID = id
	if err := h.ParseJSONBody(ctx, req); err != nil {
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
// @Tags    providers_ssl
// @Produce json
// @Id      updateSslProviderMeta
// @Param   ID path string true "provider ID"
// @Param   body body ssldto.UpdateSslMetaReq true "request data"
// @Success 200 {object} ssldto.UpdateSslMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{ID}/meta [put]
func (h *ProvidersHandler) UpdateSslMeta(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSsl,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewUpdateSslMetaReq()
	req.ID = id
	if err := h.ParseJSONBody(ctx, req); err != nil {
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
// @Tags    providers_ssl
// @Produce json
// @Id      deleteSslProvider
// @Param   ID path string true "provider ID"
// @Success 200 {object} ssldto.DeleteSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssls/{ID} [delete]
func (h *ProvidersHandler) DeleteSsl(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSsl,
		ResourceID:     id,
		Action:         base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := ssldto.NewDeleteSslReq()
	req.ID = id
	if err := h.ParseRequest(ctx, req, nil); err != nil {
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
