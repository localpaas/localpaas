package slackuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/slackuc/slackdto"
	"github.com/localpaas/localpaas/services/slack"
)

func (uc *SlackUC) TestSendSlackMsg(
	ctx context.Context,
	auth *basedto.Auth,
	req *slackdto.TestSendSlackMsgReq,
) (*slackdto.TestSendSlackMsgResp, error) {
	client := slack.NewClient()

	err := client.PostWebhook(ctx, req.Webhook, "", req.TestMsg)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &slackdto.TestSendSlackMsgResp{}, nil
}
