package systemhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/lpappuc/lpappdto"
)

// ReloadLocalPaasAppConfig Reloads LocalPaas config files
// @Summary Reloads LocalPaas config files
// @Description Reloads LocalPaas config files
// @Tags    system_localpaas_app
// @Produce json
// @Id      reloadLocalPaasAppConfig
// @Param   body body lpappdto.ReloadLpAppConfigReq true "request data"
// @Success 200 {object} lpappdto.ReloadLpAppConfigResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/localpaas/config/reload [post]
func (h *SystemHandler) ReloadLocalPaasAppConfig(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := lpappdto.NewReloadLpAppConfigReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.lpAppUC.ReloadLpAppConfig(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// RestartLocalPaasApp Restarts localpaas app containers
// @Summary Restarts localpaas app containers
// @Description Restarts localpaas app containers
// @Tags    system_localpaas_app
// @Produce json
// @Id      restartLocalPaasApp
// @Param   body body lpappdto.RestartLpAppReq true "request data"
// @Success 200 {object} lpappdto.RestartLpAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/localpaas/restart [post]
func (h *SystemHandler) RestartLocalPaasApp(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := lpappdto.NewRestartLpAppReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.lpAppUC.RestartLpApp(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
