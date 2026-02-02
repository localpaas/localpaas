package webhookhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

// HandleGithubWebhook Handles Github webhook
// @Summary Handles Github webhook
// @Description Handles Github webhook
// @Tags    webhooks
// @Produce json
// @Id      handleGithubWebhook
// @Param   body body webhookdto.HandleGitWebhookReq true "request data"
// @Success 200 {object} webhookdto.HandleGitWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /webhooks/github [post]
func (h *WebhookHandler) HandleGithubWebhook(ctx *gin.Context) {
	req := webhookdto.NewHandleGitWebhookReq()
	req.Request = ctx.Request
	req.GitSource = base.GitSourceGithub

	resp, err := h.webhookUC.HandleGitWebhook(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// HandleGitlabWebhook Handles Gitlab webhook
// @Summary Handles Gitlab webhook
// @Description Handles Gitlab webhook
// @Tags    webhooks
// @Produce json
// @Id      handleGitlabWebhook
// @Param   body body webhookdto.HandleGitWebhookReq true "request data"
// @Success 200 {object} webhookdto.HandleGitWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /webhooks/gitlab [post]
func (h *WebhookHandler) HandleGitlabWebhook(ctx *gin.Context) {
	req := webhookdto.NewHandleGitWebhookReq()
	req.Request = ctx.Request
	req.GitSource = base.GitSourceGitlab

	resp, err := h.webhookUC.HandleGitWebhook(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// HandleGiteaWebhook Handles Gitea webhook
// @Summary Handles Gitea webhook
// @Description Handles Gitea webhook
// @Tags    webhooks
// @Produce json
// @Id      handleGiteaWebhook
// @Param   body body webhookdto.HandleGitWebhookReq true "request data"
// @Success 200 {object} webhookdto.HandleGitWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /webhooks/gitea [post]
func (h *WebhookHandler) HandleGiteaWebhook(ctx *gin.Context) {
	req := webhookdto.NewHandleGitWebhookReq()
	req.Request = ctx.Request
	req.GitSource = base.GitSourceGitea

	resp, err := h.webhookUC.HandleGitWebhook(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// HandleBitbucketWebhook Handles Bitbucket webhook
// @Summary Handles Bitbucket webhook
// @Description Handles Bitbucket webhook
// @Tags    webhooks
// @Produce json
// @Id      handleBitbucketWebhook
// @Param   body body webhookdto.HandleGitWebhookReq true "request data"
// @Success 200 {object} webhookdto.HandleGitWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /webhooks/bitbucket [post]
func (h *WebhookHandler) HandleBitbucketWebhook(ctx *gin.Context) {
	req := webhookdto.NewHandleGitWebhookReq()
	req.Request = ctx.Request
	req.GitSource = base.GitSourceBitbucket

	resp, err := h.webhookUC.HandleGitWebhook(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
