package discorduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
)

func (uc *DiscordUC) ListDiscord(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.ListDiscordReq,
) (*discorddto.ListDiscordResp, error) {
	req.Type = currentSettingType
	resp, err := providers.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &providers.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := discorddto.TransformDiscords(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &discorddto.ListDiscordResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
