package slackuc

import (
	"context"
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/slackuc/slackdto"
)

func (uc *SlackUC) UpdateSlack(
	ctx context.Context,
	auth *basedto.Auth,
	req *slackdto.UpdateSlackReq,
) (*slackdto.UpdateSlackResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		slackData := &updateSlackData{}
		err := uc.loadSlackDataForUpdate(ctx, db, req, slackData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingSlackData{}
		uc.prepareUpdatingSlack(req.SlackBaseReq, slackData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &slackdto.UpdateSlackResp{}, nil
}

type updateSlackData struct {
	Setting *entity.Setting
}

func (uc *SlackUC) loadSlackDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *slackdto.UpdateSlackReq,
	data *updateSlackData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeSlack, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if req.Name != "" && !strings.EqualFold(setting.Name, req.Name) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeSlack, req.Name, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("Slack").
				WithMsgLog("slack setting '%s' already exists", req.Name)
		}
	}

	return nil
}

func (uc *SlackUC) prepareUpdatingSlack(
	req *slackdto.SlackBaseReq,
	data *updateSlackData,
	persistingData *persistingSlackData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.UpdateVer++
	setting.UpdatedAt = timeNow
	setting.Name = gofn.Coalesce(req.Name, setting.Name)

	slack := &entity.Slack{
		Webhook: entity.NewEncryptedField(req.Webhook),
	}
	setting.MustSetData(slack)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
