package systemhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/nginxuc/nginxdto"
)

// ReloadNginxConfig Reloads nginx config files
// @Summary Reloads nginx config files
// @Description Reloads nginx config files
// @Tags    system_nginx
// @Produce json
// @Id      reloadNginxConfig
// @Param   body body nginxdto.ReloadNginxConfigReq true "request data"
// @Success 200 {object} nginxdto.ReloadNginxConfigResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/nginx/config/reload [post]
func (h *SystemHandler) ReloadNginxConfig(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := nginxdto.NewReloadNginxConfigReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.nginxUC.ReloadNginxConfig(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ResetNginxConfig Resets nginx config files
// @Summary Resets nginx config files
// @Description Resets nginx config files
// @Tags    system_nginx
// @Produce json
// @Id      resetNginxConfig
// @Param   body body nginxdto.ResetNginxConfigReq true "request data"
// @Success 200 {object} nginxdto.ResetNginxConfigResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/nginx/config/reset [post]
func (h *SystemHandler) ResetNginxConfig(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := nginxdto.NewResetNginxConfigReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.nginxUC.ResetNginxConfig(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// RestartNginx Restarts nginx containers
// @Summary Restarts nginx containers
// @Description Restarts nginx containers
// @Tags    system_nginx
// @Produce json
// @Id      restartNginx
// @Param   body body nginxdto.RestartNginxReq true "request data"
// @Success 200 {object} nginxdto.RestartNginxResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/nginx/restart [post]
func (h *SystemHandler) RestartNginx(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := nginxdto.NewRestartNginxReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.nginxUC.RestartNginx(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
