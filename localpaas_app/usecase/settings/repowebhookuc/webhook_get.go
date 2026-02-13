package repowebhookuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/repowebhookuc/repowebhookdto"
)

func (uc *RepoWebhookUC) GetRepoWebhook(
	ctx context.Context,
	auth *basedto.Auth,
	req *repowebhookdto.GetRepoWebhookReq,
) (*repowebhookdto.GetRepoWebhookResp, error) {
	req.Type = currentSettingType
	setting, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsRepoWebhook().MustDecrypt()
	resp, err := repowebhookdto.TransformRepoWebhook(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &repowebhookdto.GetRepoWebhookResp{
		Data: resp,
	}, nil
}
