package slackuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc/slackdto"
)

func (uc *SlackUC) CreateSlack(
	ctx context.Context,
	auth *basedto.Auth,
	req *slackdto.CreateSlackReq,
) (*slackdto.CreateSlackResp, error) {
	slackData := &createSlackData{}
	err := uc.loadSlackData(ctx, uc.db, req, slackData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingSlackData{}
	uc.preparePersistingSlack(req.SlackBaseReq, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &slackdto.CreateSlackResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createSlackData struct {
}

func (uc *SlackUC) loadSlackData(
	ctx context.Context,
	db database.IDB,
	req *slackdto.CreateSlackReq,
	_ *createSlackData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeSlack, req.Name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("Slack").
			WithMsgLog("slack setting '%s' already exists", req.Name)
	}

	return nil
}

type persistingSlackData struct {
	settingservice.PersistingSettingData
}

func (uc *SlackUC) preparePersistingSlack(
	req *slackdto.SlackBaseReq,
	persistingData *persistingSlackData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeSlack,
		Status:    base.SettingStatusActive,
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	slack := &entity.Slack{
		Webhook: req.Webhook,
	}
	setting.MustSetData(slack)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *SlackUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingSlackData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
