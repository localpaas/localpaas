package discorduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
	"github.com/localpaas/localpaas/services/discord"
)

func (uc *DiscordUC) TestSendDiscordMsg(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.TestSendDiscordMsgReq,
) (*discorddto.TestSendDiscordMsgResp, error) {
	client := discord.NewClient()

	_, err := client.WebhookExecute(ctx, req.Webhook, true, req.TestMsg)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &discorddto.TestSendDiscordMsgResp{}, nil
}
