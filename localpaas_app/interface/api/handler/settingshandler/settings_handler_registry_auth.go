package settingshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListRegistryAuth Lists registry auth settings
// @Summary Lists registry auth settings
// @Description Lists registry auth settings
// @Tags    settings_registry_auth
// @Produce json
// @Id      listRegistryAuthSettings
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} registryauthdto.ListRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/registry-auth [get]
func (h *SettingsHandler) ListRegistryAuth(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		ResourceType:   base.ResourceTypeRegistryAuth,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewListRegistryAuthReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
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

// GetRegistryAuth Gets registry auth setting details
// @Summary Gets registry auth setting details
// @Description Gets registry auth setting details
// @Tags    settings_registry_auth
// @Produce json
// @Id      getRegistryAuthSetting
// @Param   ID path string true "setting ID"
// @Success 200 {object} registryauthdto.GetRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/registry-auth/{ID} [get]
func (h *SettingsHandler) GetRegistryAuth(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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
	if err = h.ParseRequest(ctx, req, nil); err != nil {
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

// CreateRegistryAuth Creates a new registry auth setting
// @Summary Creates a new registry auth setting
// @Description Creates a new registry auth setting
// @Tags    settings_registry_auth
// @Produce json
// @Id      createRegistryAuthSetting
// @Param   body body registryauthdto.CreateRegistryAuthReq true "request data"
// @Success 201 {object} registryauthdto.CreateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/registry-auth [post]
func (h *SettingsHandler) CreateRegistryAuth(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := registryauthdto.NewCreateRegistryAuthReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
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
// @Tags    settings_registry_auth
// @Produce json
// @Id      updateRegistryAuthSetting
// @Param   ID path string true "setting ID"
// @Param   body body registryauthdto.UpdateRegistryAuthReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/registry-auth/{ID} [put]
func (h *SettingsHandler) UpdateRegistryAuth(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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
	if err := h.ParseJSONBody(ctx, req); err != nil {
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
// @Tags    settings_registry_auth
// @Produce json
// @Id      updateRegistryAuthMetaSetting
// @Param   ID path string true "setting ID"
// @Param   body body registryauthdto.UpdateRegistryAuthMetaReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/registry-auth/{ID}/meta [put]
func (h *SettingsHandler) UpdateRegistryAuthMeta(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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
	if err := h.ParseJSONBody(ctx, req); err != nil {
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

// DeleteRegistryAuth Deletes registry auth setting
// @Summary Deletes registry auth setting
// @Description Deletes registry auth setting
// @Tags    settings_registry_auth
// @Produce json
// @Id      deleteRegistryAuthSetting
// @Param   ID path string true "setting ID"
// @Success 200 {object} registryauthdto.DeleteRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/registry-auth/{ID} [delete]
func (h *SettingsHandler) DeleteRegistryAuth(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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
	if err := h.ParseRequest(ctx, req, nil); err != nil {
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
