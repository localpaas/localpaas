package basicauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
)

func (uc *BasicAuthUC) UpdateBasicAuthMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.UpdateBasicAuthMetaReq,
) (*basicauthdto.UpdateBasicAuthMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		authData := &updateBasicAuthData{}
		err := uc.loadBasicAuthDataForUpdateMeta(ctx, db, req, authData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingBasicAuthMeta(req, authData)
		return uc.persistBasicAuthMeta(ctx, db, authData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &basicauthdto.UpdateBasicAuthMetaResp{}, nil
}

func (uc *BasicAuthUC) loadBasicAuthDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *basicauthdto.UpdateBasicAuthMetaReq,
	data *updateBasicAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeBasicAuth, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	return nil
}

func (uc *BasicAuthUC) prepareUpdatingBasicAuthMeta(
	req *basicauthdto.UpdateBasicAuthMetaReq,
	data *updateBasicAuthData,
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

func (uc *BasicAuthUC) persistBasicAuthMeta(
	ctx context.Context,
	db database.IDB,
	data *updateBasicAuthData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
