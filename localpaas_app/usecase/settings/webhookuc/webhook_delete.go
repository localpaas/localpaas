package webhookuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/webhookuc/webhookdto"
)

func (uc *WebhookUC) DeleteWebhook(
	ctx context.Context,
	auth *basedto.Auth,
	req *webhookdto.DeleteWebhookReq,
) (*webhookdto.DeleteWebhookResp, error) {
	req.Type = currentSettingType
	_, err := settings.DeleteSetting(ctx, uc.db, &req.DeleteSettingReq, &settings.DeleteSettingData{
		SettingRepo:              uc.settingRepo,
		ProjectSharedSettingRepo: uc.projectSharedSettingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &webhookdto.DeleteWebhookResp{}, nil
}
