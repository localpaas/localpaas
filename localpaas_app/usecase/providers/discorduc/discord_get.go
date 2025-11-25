package discorduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
)

func (uc *DiscordUC) GetDiscord(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.GetDiscordReq,
) (*discorddto.GetDiscordResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, base.SettingTypeDiscord, req.ID, false)
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
