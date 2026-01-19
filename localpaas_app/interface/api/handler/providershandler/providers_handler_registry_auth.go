package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc/registryauthdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListRegistryAuth Lists registry auth providers
// @Summary Lists registry auth providers
// @Description Lists registry auth providers
// @Tags    providers_registry_auth
// @Produce json
// @Id      listRegistryAuthProviders
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} registryauthdto.ListRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth [get]
func (h *ProvidersHandler) ListRegistryAuth(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeRegistryAuth,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewListRegistryAuthReq()
	req.GlobalOnly = true
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
// @Tags    providers_registry_auth
// @Produce json
// @Id      getRegistryAuthProvider
// @Param   id path string true "provider ID"
// @Success 200 {object} registryauthdto.GetRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/{id} [get]
func (h *ProvidersHandler) GetRegistryAuth(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "id")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeRegistryAuth,
		ResourceID:     id,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewGetRegistryAuthReq()
	req.ID = id
	req.GlobalOnly = true
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
// @Tags    providers_registry_auth
// @Produce json
// @Id      createRegistryAuthProvider
// @Param   body body registryauthdto.CreateRegistryAuthReq true "request data"
// @Success 201 {object} registryauthdto.CreateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth [post]
func (h *ProvidersHandler) CreateRegistryAuth(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewCreateRegistryAuthReq()
	req.GlobalOnly = true
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
// @Tags    providers_registry_auth
// @Produce json
// @Id      updateRegistryAuthProvider
// @Param   id path string true "provider ID"
// @Param   body body registryauthdto.UpdateRegistryAuthReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/{id} [put]
func (h *ProvidersHandler) UpdateRegistryAuth(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "id")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeRegistryAuth,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewUpdateRegistryAuthReq()
	req.ID = id
	req.GlobalOnly = true
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
// @Tags    providers_registry_auth
// @Produce json
// @Id      updateRegistryAuthProviderMeta
// @Param   id path string true "provider ID"
// @Param   body body registryauthdto.UpdateRegistryAuthMetaReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/{id}/meta [put]
func (h *ProvidersHandler) UpdateRegistryAuthMeta(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "id")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeRegistryAuth,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewUpdateRegistryAuthMetaReq()
	req.ID = id
	req.GlobalOnly = true
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
// @Tags    providers_registry_auth
// @Produce json
// @Id      deleteRegistryAuthProvider
// @Param   id path string true "provider ID"
// @Success 200 {object} registryauthdto.DeleteRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/{id} [delete]
func (h *ProvidersHandler) DeleteRegistryAuth(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "id")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeRegistryAuth,
		ResourceID:     id,
		Action:         base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewDeleteRegistryAuthReq()
	req.ID = id
	req.GlobalOnly = true
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

// TestRegistryAuthConn Tests registry auth connection
// @Summary Tests registry auth connection
// @Description Tests registry auth connection
// @Tags    providers_registry_auth
// @Produce json
// @Id      testRegistryAuthConn
// @Param   body body registryauthdto.TestRegistryAuthConnReq true "request data"
// @Success 200 {object} registryauthdto.TestRegistryAuthConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/registry-auth/test-conn [post]
func (h *ProvidersHandler) TestRegistryAuthConn(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewTestRegistryAuthConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.registryAuthUC.TestRegistryAuthConn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
