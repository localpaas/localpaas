package webhookhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

// HandleWebhookGithub Handles Github webhook
// @Summary Handles Github webhook
// @Description Handles Github webhook
// @Tags    webhooks
// @Produce json
// @Id      handleWebhookGithub
// @Param   body body webhookdto.HandleWebhookGithubReq true "request data"
// @Success 200 {object} webhookdto.HandleWebhookGithubResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /webhooks/github [post]
func (h *WebhookHandler) HandleWebhookGithub(ctx *gin.Context) {
	req := webhookdto.NewHandleWebhookGithubReq()
	req.Request = ctx.Request

	resp, err := h.webhookUC.HandleWebhookGithub(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
