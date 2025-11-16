package settingshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListSsl Lists SSL settings
// @Summary Lists SSL settings
// @Description Lists SSL settings
// @Tags    settings_ssl
// @Produce json
// @Id      listSslSettings
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} ssldto.ListSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls [get]
func (h *SettingsHandler) ListSsl(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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

// GetSsl Gets SSL setting details
// @Summary Gets SSL setting details
// @Description Gets SSL setting details
// @Tags    settings_ssl
// @Produce json
// @Id      getSslSetting
// @Param   ID path string true "setting ID"
// @Success 200 {object} ssldto.GetSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls/{ID} [get]
func (h *SettingsHandler) GetSsl(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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

// CreateSsl Creates a new SSL setting
// @Summary Creates a new SSL setting
// @Description Creates a new SSL setting
// @Tags    settings_ssl
// @Produce json
// @Id      createSslSetting
// @Param   body body ssldto.CreateSslReq true "request data"
// @Success 201 {object} ssldto.CreateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls [post]
func (h *SettingsHandler) CreateSsl(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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
// @Tags    settings_ssl
// @Produce json
// @Id      updateSslSetting
// @Param   ID path string true "setting ID"
// @Success 200 {object} ssldto.UpdateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls/{ID} [put]
func (h *SettingsHandler) UpdateSsl(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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

// DeleteSsl Deletes SSL setting
// @Summary Deletes SSL setting
// @Description Deletes SSL setting
// @Tags    settings_ssl
// @Produce json
// @Id      deleteSslSetting
// @Param   ID path string true "setting ID"
// @Success 200 {object} ssldto.DeleteSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssls/{ID} [delete]
func (h *SettingsHandler) DeleteSsl(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
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
