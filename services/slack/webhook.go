package slack

import (
	"context"

	"github.com/slack-go/slack"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type WebhookMessageOption func(*slack.WebhookMessage)

func (c *Client) PostWebhook(ctx context.Context, webhookURL, channel, text string,
	options ...WebhookMessageOption) error {
	msg := &slack.WebhookMessage{
		Channel: channel,
		Text:    text,
	}
	for _, opt := range options {
		opt(msg)
	}
	err := slack.PostWebhookCustomHTTPContext(ctx, webhookURL, c.getHttpClient(), msg)
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}
