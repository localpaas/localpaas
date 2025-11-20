package basicauthuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
)

func (uc *BasicAuthUC) DeleteBasicAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.DeleteBasicAuthReq,
) (*basicauthdto.DeleteBasicAuthResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		basicAuthData := &deleteBasicAuthData{}
		err := uc.loadBasicAuthDataForDelete(ctx, db, req, basicAuthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingBasicAuthData{}
		uc.prepareDeletingBasicAuth(basicAuthData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &basicauthdto.DeleteBasicAuthResp{}, nil
}

type deleteBasicAuthData struct {
	Setting *entity.Setting
}

func (uc *BasicAuthUC) loadBasicAuthDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *basicauthdto.DeleteBasicAuthReq,
	data *deleteBasicAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeBasicAuth, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *BasicAuthUC) prepareDeletingBasicAuth(
	data *deleteBasicAuthData,
	persistingData *persistingBasicAuthData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
