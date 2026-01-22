package slackuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc/slackdto"
)

func (uc *SlackUC) ListSlack(
	ctx context.Context,
	auth *basedto.Auth,
	req *slackdto.ListSlackReq,
) (*slackdto.ListSlackResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := slackdto.TransformSlacks(resp.Data, req.ObjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &slackdto.ListSlackResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
