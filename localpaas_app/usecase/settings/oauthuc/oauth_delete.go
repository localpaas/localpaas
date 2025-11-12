package oauthuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
)

func (uc *OAuthUC) DeleteOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.DeleteOAuthReq,
) (*oauthdto.DeleteOAuthResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		oauthData := &deleteOAuthData{}
		err := uc.loadOAuthDataForDelete(ctx, db, req, oauthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingOAuthData{}
		uc.prepareDeletingOAuth(oauthData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &oauthdto.DeleteOAuthResp{}, nil
}

type deleteOAuthData struct {
	Setting *entity.Setting
}

func (uc *OAuthUC) loadOAuthDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *oauthdto.DeleteOAuthReq,
	data *deleteOAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, req.ID,
		bunex.SelectFor("UPDATE OF setting"),
		bunex.SelectWhere("setting.type = ?", base.SettingTypeOAuth),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *OAuthUC) prepareDeletingOAuth(
	data *deleteOAuthData,
	persistingData *persistingOAuthData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
