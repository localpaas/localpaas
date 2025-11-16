package apikeyuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) UpdateAPIKeyMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.UpdateAPIKeyMetaReq,
) (*apikeydto.UpdateAPIKeyMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		apikeyData := &updateAPIKeyData{}
		err := uc.loadAPIKeyDataForUpdateMeta(ctx, db, req, apikeyData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingAPIKeyMeta(req, apikeyData)
		return uc.persistAPIKeyMeta(ctx, db, apikeyData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.UpdateAPIKeyMetaResp{}, nil
}

type updateAPIKeyData struct {
	Setting *entity.Setting
}

func (uc *APIKeyUC) loadAPIKeyDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *apikeydto.UpdateAPIKeyMetaReq,
	data *updateAPIKeyData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeAPIKey, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *APIKeyUC) prepareUpdatingAPIKeyMeta(
	req *apikeydto.UpdateAPIKeyMetaReq,
	data *updateAPIKeyData,
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

func (uc *APIKeyUC) persistAPIKeyMeta(
	ctx context.Context,
	db database.IDB,
	data *updateAPIKeyData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
