package webhookuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/webhookuc/webhookdto"
)

const (
	currentSettingType    = base.SettingTypeWebhook
	currentSettingVersion = entity.CurrentWebhookVersion
)

func (uc *WebhookUC) CreateWebhook(
	ctx context.Context,
	auth *basedto.Auth,
	req *webhookdto.CreateWebhookReq,
) (*webhookdto.CreateWebhookResp, error) {
	req.Type = currentSettingType
	resp, err := settings.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &settings.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		Version:       currentSettingVersion,
		PrepareCreation: func(
			ctx context.Context,
			db database.Tx,
			data *settings.CreateSettingData,
			pData *settings.PersistingSettingCreationData,
		) error {
			pData.Setting.Kind = string(req.Kind)
			err := pData.Setting.SetData(req.ToEntity())
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &webhookdto.CreateWebhookResp{
		Data: resp.Data,
	}, nil
}
