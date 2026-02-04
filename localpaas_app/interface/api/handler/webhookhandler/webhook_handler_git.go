package webhookhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

// HandleRepoWebhook Handles Repo webhook
// @Summary Handles Repo webhook
// @Description Handles Repo webhook
// @Tags    webhooks
// @Produce json
// @Id      handleRepoWebhook
// @Param   kind path string true "webhook kind"
// @Param   secret path string true "webhook secret"
// @Param   body body webhookdto.HandleRepoWebhookReq true "request data"
// @Success 200 {object} webhookdto.HandleRepoWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /webhooks/{kind}/{secret} [post]
func (h *WebhookHandler) HandleRepoWebhook(ctx *gin.Context) {
	kind, err := h.ParseStringParam(ctx, "kind")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	secret, err := h.ParseStringParam(ctx, "secret")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := webhookdto.NewHandleRepoWebhookReq()
	req.Request = ctx.Request
	req.WebhookKind = base.WebhookKind(kind)
	req.Secret = secret

	resp, err := h.webhookUC.HandleRepoWebhook(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
