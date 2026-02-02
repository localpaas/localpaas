package webhookuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/webhookuc/webhookdto"
)

func (uc *WebhookUC) ListWebhook(
	ctx context.Context,
	auth *basedto.Auth,
	req *webhookdto.ListWebhookReq,
) (*webhookdto.ListWebhookResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := webhookdto.TransformWebhooks(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &webhookdto.ListWebhookResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
