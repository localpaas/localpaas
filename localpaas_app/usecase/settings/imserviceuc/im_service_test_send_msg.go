package imserviceuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc/imservicedto"
	"github.com/localpaas/localpaas/services/discord"
	"github.com/localpaas/localpaas/services/slack"
)

func (uc *IMServiceUC) TestSendInstantMsg(
	ctx context.Context,
	auth *basedto.Auth,
	req *imservicedto.TestSendInstantMsgReq,
) (*imservicedto.TestSendInstantMsgResp, error) {
	switch {
	case req.Slack != nil:
		if err := uc.testSendSlackMsg(ctx, req); err != nil {
			return nil, apperrors.Wrap(err)
		}
	case req.Discord != nil:
		if err := uc.testSendDiscordMsg(ctx, req); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	return &imservicedto.TestSendInstantMsgResp{}, nil
}

func (uc *IMServiceUC) testSendSlackMsg(
	ctx context.Context,
	req *imservicedto.TestSendInstantMsgReq,
) error {
	client := slack.NewClient()

	err := client.PostWebhook(ctx, req.Slack.Webhook, "", req.TestMsg)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *IMServiceUC) testSendDiscordMsg(
	ctx context.Context,
	req *imservicedto.TestSendInstantMsgReq,
) error {
	client := discord.NewClient()

	_, err := client.WebhookExecute(ctx, req.Discord.Webhook, true, req.TestMsg)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
