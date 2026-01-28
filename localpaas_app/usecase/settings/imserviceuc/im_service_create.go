package imserviceuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc/imservicedto"
)

const (
	currentSettingType    = base.SettingTypeIMService
	currentSettingVersion = entity.CurrentIMServiceVersion
)

func (uc *IMServiceUC) CreateIMService(
	ctx context.Context,
	auth *basedto.Auth,
	req *imservicedto.CreateIMServiceReq,
) (*imservicedto.CreateIMServiceResp, error) {
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
			imService := &entity.IMService{}
			switch {
			case req.Slack != nil:
				pData.Setting.Kind = string(base.IMServiceKindSlack)
				imService.Slack = &entity.Slack{
					Webhook: entity.NewEncryptedField(req.Slack.Webhook),
				}

			case req.Discord != nil:
				pData.Setting.Kind = string(base.IMServiceKindDiscord)
				imService.Discord = &entity.Discord{
					Webhook: entity.NewEncryptedField(req.Discord.Webhook),
				}
			}

			err := pData.Setting.SetData(imService)
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &imservicedto.CreateIMServiceResp{
		Data: resp.Data,
	}, nil
}
