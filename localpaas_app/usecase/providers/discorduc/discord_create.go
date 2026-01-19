package discorduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
)

const (
	currentSettingType    = base.SettingTypeDiscord
	currentSettingVersion = entity.CurrentDiscordVersion
)

func (uc *DiscordUC) CreateDiscord(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.CreateDiscordReq,
) (*discorddto.CreateDiscordResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		Version:       currentSettingVersion,
		PrepareCreation: func(ctx context.Context, db database.Tx, data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData) error {
			err := pData.Setting.SetData(&entity.Discord{
				Webhook: entity.NewEncryptedField(req.Webhook),
			})
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &discorddto.CreateDiscordResp{
		Data: resp.Data,
	}, nil
}
