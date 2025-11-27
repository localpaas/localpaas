package discorduc

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
)

func (uc *DiscordUC) CreateDiscord(
	ctx context.Context,
	auth *basedto.Auth,
	req *discorddto.CreateDiscordReq,
) (*discorddto.CreateDiscordResp, error) {
	discordData := &createDiscordData{}
	err := uc.loadDiscordData(ctx, uc.db, req, discordData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingDiscordData{}
	uc.preparePersistingDiscord(req.DiscordBaseReq, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &discorddto.CreateDiscordResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createDiscordData struct {
}

func (uc *DiscordUC) loadDiscordData(
	ctx context.Context,
	db database.IDB,
	req *discorddto.CreateDiscordReq,
	_ *createDiscordData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeDiscord, req.Name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("Discord").
			WithMsgLog("discord setting '%s' already exists", req.Name)
	}

	return nil
}

type persistingDiscordData struct {
	settingservice.PersistingSettingData
}

func (uc *DiscordUC) preparePersistingDiscord(
	req *discorddto.DiscordBaseReq,
	persistingData *persistingDiscordData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeDiscord,
		Status:    base.SettingStatusActive,
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	discord := &entity.Discord{
		Webhook: entity.NewEncryptedField(req.Webhook),
	}
	setting.MustSetData(discord)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *DiscordUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingDiscordData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
