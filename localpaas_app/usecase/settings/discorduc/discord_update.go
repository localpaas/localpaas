package discorduc

import (
	"context"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/discorduc/discorddto"
)

func (uc *DiscordUC) UpdateDiscord(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.UpdateDiscordReq,
) (*discorddto.UpdateDiscordResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		discordData := &updateDiscordData{}
		err := uc.loadDiscordDataForUpdate(ctx, db, req, discordData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingDiscordData{}
		uc.prepareUpdatingDiscord(req.DiscordBaseReq, discordData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &discorddto.UpdateDiscordResp{}, nil
}

type updateDiscordData struct {
	Setting *entity.Setting
}

func (uc *DiscordUC) loadDiscordDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *discorddto.UpdateDiscordReq,
	data *updateDiscordData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeDiscord, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if req.Name != "" && !strings.EqualFold(setting.Name, req.Name) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeDiscord, req.Name, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("Discord").
				WithMsgLog("discord setting '%s' already exists", req.Name)
		}
	}

	return nil
}

func (uc *DiscordUC) prepareUpdatingDiscord(
	req *discorddto.DiscordBaseReq,
	data *updateDiscordData,
	persistingData *persistingDiscordData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting

	if req.Name != "" {
		setting.Name = req.Name
	}
	discord := &entity.Discord{
		Webhook: req.Webhook,
	}
	setting.MustSetData(discord.MustEncrypt())

	setting.UpdatedAt = timeNow
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
