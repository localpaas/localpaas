package webhookhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

// HandleGitWebhook Handles Git webhook
// @Summary Handles Git webhook
// @Description Handles Git webhook
// @Tags    webhooks
// @Produce json
// @Id      handleGitWebhook
// @Param   gitSource path string true "git source"
// @Param   secret path string true "webhook secret"
// @Param   body body webhookdto.HandleGitWebhookReq true "request data"
// @Success 200 {object} webhookdto.HandleGitWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /webhooks/{gitSource}/{secret} [post]
func (h *WebhookHandler) HandleGitWebhook(ctx *gin.Context) {
	gitSource, err := h.ParseStringParam(ctx, "gitSource")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	secret, err := h.ParseStringParam(ctx, "secret")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := webhookdto.NewHandleGitWebhookReq()
	req.Request = ctx.Request
	req.GitSource = base.GitSource(gitSource)
	req.Secret = secret

	resp, err := h.webhookUC.HandleGitWebhook(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
