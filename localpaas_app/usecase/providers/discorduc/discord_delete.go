package discorduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
)

func (uc *DiscordUC) DeleteDiscord(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.DeleteDiscordReq,
) (*discorddto.DeleteDiscordResp, error) {
	req.Type = currentSettingType
	_, err := providers.DeleteSetting(ctx, uc.db, &req.DeleteSettingReq, &providers.DeleteSettingData{
		SettingRepo:              uc.settingRepo,
		ProjectSharedSettingRepo: uc.projectSharedSettingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &discorddto.DeleteDiscordResp{}, nil
}
