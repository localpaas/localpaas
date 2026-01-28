package imserviceuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc/imservicedto"
)

func (uc *IMServiceUC) UpdateIMService(
	ctx context.Context,
	auth *basedto.Auth,
	req *imservicedto.UpdateIMServiceReq,
) (*imservicedto.UpdateIMServiceResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &settings.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *settings.UpdateSettingData,
			pData *settings.PersistingSettingData,
		) error {
			pData.Setting.Name = gofn.Coalesce(req.Name, pData.Setting.Name)

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

	return &imservicedto.UpdateIMServiceResp{}, nil
}
