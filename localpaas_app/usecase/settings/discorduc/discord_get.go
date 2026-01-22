package discorduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/discorduc/discorddto"
)

func (uc *DiscordUC) GetDiscord(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.GetDiscordReq,
) (*discorddto.GetDiscordResp, error) {
	req.Type = currentSettingType
	setting, err := settings.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &settings.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsDiscord().MustDecrypt()
	resp, err := discorddto.TransformDiscord(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &discorddto.GetDiscordResp{
		Data: resp,
	}, nil
}
