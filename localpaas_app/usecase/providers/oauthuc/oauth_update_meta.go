package oauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc/oauthdto"
)

func (uc *OAuthUC) UpdateOAuthMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.UpdateOAuthMetaReq,
) (*oauthdto.UpdateOAuthMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		oauthData := &updateOAuthData{}
		err := uc.loadOAuthDataForUpdateMeta(ctx, db, req, oauthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingOAuthMeta(req, oauthData)
		return uc.persistOAuthMeta(ctx, db, oauthData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &oauthdto.UpdateOAuthMetaResp{}, nil
}

func (uc *OAuthUC) loadOAuthDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *oauthdto.UpdateOAuthMetaReq,
	data *updateOAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeOAuth, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *OAuthUC) prepareUpdatingOAuthMeta(
	req *oauthdto.UpdateOAuthMetaReq,
	data *updateOAuthData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting

	if req.Status != nil {
		setting.Status = *req.Status
	}
	if req.ExpireAt != nil {
		setting.ExpireAt = *req.ExpireAt
	}

	setting.UpdatedAt = timeNow
}

func (uc *OAuthUC) persistOAuthMeta(
	ctx context.Context,
	db database.IDB,
	data *updateOAuthData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
