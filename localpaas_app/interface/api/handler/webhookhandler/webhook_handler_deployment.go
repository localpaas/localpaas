package webhookhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

// WebhookDeployApp Deploys an app
// @Summary Deploys an app
// @Description Deploys an app
// @Tags    webhooks
// @Produce json
// @Id      webhookDeployApp
// @Param   appToken path string true "app token"
// @Param   body body webhookdto.DeployAppReq true "request data"
// @Success 200 {object} webhookdto.DeployAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /webhooks/apps/{appToken}/deploy [post]
func (h *WebhookHandler) WebhookDeployApp(ctx *gin.Context) {
	appToken, err := h.ParseStringParam(ctx, "appToken")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := webhookdto.NewDeployAppReq()
	req.AppToken = appToken
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.webhookUC.DeployApp(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
