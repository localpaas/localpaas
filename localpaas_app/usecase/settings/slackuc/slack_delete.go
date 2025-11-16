package slackuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc/slackdto"
)

func (uc *SlackUC) DeleteSlack(
	ctx context.Context,
	auth *basedto.Auth,
	req *slackdto.DeleteSlackReq,
) (*slackdto.DeleteSlackResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		slackData := &deleteSlackData{}
		err := uc.loadSlackDataForDelete(ctx, db, req, slackData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingSlackData{}
		uc.prepareDeletingSlack(slackData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &slackdto.DeleteSlackResp{}, nil
}

type deleteSlackData struct {
	Setting *entity.Setting
}

func (uc *SlackUC) loadSlackDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *slackdto.DeleteSlackReq,
	data *deleteSlackData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeSlack, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *SlackUC) prepareDeletingSlack(
	data *deleteSlackData,
	persistingData *persistingSlackData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
