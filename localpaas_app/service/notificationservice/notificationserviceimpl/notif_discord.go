package notificationserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/services/im/discord"
)

func (s *service) discordSendMsg(
	ctx context.Context,
	setting *entity.IMDiscord,
	msg string,
) error {
	webhookURL, err := setting.Webhook.GetPlain()
	if err != nil {
		return apperrors.New(err)
	}
	_, err = discord.NewClient().WebhookExecute(ctx, webhookURL, true, msg)
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}
