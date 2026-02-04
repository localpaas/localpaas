package repowebhookuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/repowebhookuc/repowebhookdto"
)

func (uc *RepoWebhookUC) UpdateRepoWebhookMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *repowebhookdto.UpdateRepoWebhookMetaReq,
) (*repowebhookdto.UpdateRepoWebhookMetaResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &settings.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &repowebhookdto.UpdateRepoWebhookMetaResp{}, nil
}
