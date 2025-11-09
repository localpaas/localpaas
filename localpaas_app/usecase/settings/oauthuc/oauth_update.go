package oauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

func (uc *OAuthUC) UpdateOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.UpdateOAuthReq,
) (*oauthdto.UpdateOAuthResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		oauthData := &updateOAuthData{}
		err := uc.loadOAuthDataForUpdate(ctx, db, req, oauthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingOAuthData{}
		err = uc.prepareUpdatingOAuth(req.OAuthBaseReq, oauthData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &oauthdto.UpdateOAuthResp{}, nil
}

type updateOAuthData struct {
	Setting *entity.Setting
}

func (uc *OAuthUC) loadOAuthDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *oauthdto.UpdateOAuthReq,
	data *updateOAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, req.ID,
		bunex.SelectFor("UPDATE OF setting"),
		bunex.SelectWhere("setting.type = ?", base.SettingTypeOAuth),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	uc.preprocessRequest(base.OAuthType(setting.Name), req.OAuthBaseReq)
	return nil
}

func (uc *OAuthUC) prepareUpdatingOAuth(
	req *oauthdto.OAuthBaseReq,
	data *updateOAuthData,
	persistingData *persistingOAuthData,
) (err error) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting

	oauth := &entity.OAuth{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		Organization: req.Organization,
		CallbackURL:  req.CallbackURL,
		AuthURL:      req.AuthURL,
		TokenURL:     req.TokenURL,
		ProfileURL:   req.ProfileURL,
		Scopes:       req.Scopes,
	}
	err = setting.SetData(oauth)
	if err != nil {
		return apperrors.Wrap(err)
	}

	setting.UpdatedAt = timeNow
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	return nil
}
