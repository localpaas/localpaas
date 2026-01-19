package slackuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/slackuc/slackdto"
)

func (uc *SlackUC) ListSlack(
	ctx context.Context,
	auth *basedto.Auth,
	req *slackdto.ListSlackReq,
) (*slackdto.ListSlackResp, error) {
	req.Type = currentSettingType
	resp, err := providers.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &providers.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := slackdto.TransformSlacks(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &slackdto.ListSlackResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
