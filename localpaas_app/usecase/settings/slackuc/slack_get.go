package slackuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc/slackdto"
)

func (uc *SlackUC) GetSlack(
	ctx context.Context,
	auth *basedto.Auth,
	req *slackdto.GetSlackReq,
) (*slackdto.GetSlackResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, req.ID,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeSlack),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := slackdto.TransformSlack(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &slackdto.GetSlackResp{
		Data: resp,
	}, nil
}
