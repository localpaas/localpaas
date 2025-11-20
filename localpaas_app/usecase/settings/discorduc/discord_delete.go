package discorduc

import (
	"context"

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

func (uc *DiscordUC) DeleteDiscord(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.DeleteDiscordReq,
) (*discorddto.DeleteDiscordResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		discordData := &deleteDiscordData{}
		err := uc.loadDiscordDataForDelete(ctx, db, req, discordData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingDiscordData{}
		uc.prepareDeletingDiscord(discordData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &discorddto.DeleteDiscordResp{}, nil
}

type deleteDiscordData struct {
	Setting *entity.Setting
}

func (uc *DiscordUC) loadDiscordDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *discorddto.DeleteDiscordReq,
	data *deleteDiscordData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeDiscord, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *DiscordUC) prepareDeletingDiscord(
	data *deleteDiscordData,
	persistingData *persistingDiscordData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
