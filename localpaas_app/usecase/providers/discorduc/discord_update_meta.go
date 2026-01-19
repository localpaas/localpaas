package discorduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
)

func (uc *DiscordUC) UpdateDiscordMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.UpdateDiscordMetaReq,
) (*discorddto.UpdateDiscordMetaResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &providers.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &discorddto.UpdateDiscordMetaResp{}, nil
}
