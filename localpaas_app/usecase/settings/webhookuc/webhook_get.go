package webhookuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/webhookuc/webhookdto"
)

func (uc *WebhookUC) GetWebhook(
	ctx context.Context,
	auth *basedto.Auth,
	req *webhookdto.GetWebhookReq,
) (*webhookdto.GetWebhookResp, error) {
	req.Type = currentSettingType
	setting, err := settings.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &settings.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsWebhook().MustDecrypt()
	resp, err := webhookdto.TransformWebhook(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &webhookdto.GetWebhookResp{
		Data: resp,
	}, nil
}
