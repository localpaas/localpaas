package secretuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
)

func (uc *SecretUC) UpdateSecretMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *secretdto.UpdateSecretMetaReq,
) (*secretdto.UpdateSecretMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		secretData := &updateSecretData{}
		err := uc.loadSecretDataForUpdateMeta(ctx, db, req, secretData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingSecretMeta(req, secretData)
		return uc.persistSecretMeta(ctx, db, secretData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &secretdto.UpdateSecretMetaResp{}, nil
}

type updateSecretData struct {
	Setting *entity.Setting
}

func (uc *SecretUC) loadSecretDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *secretdto.UpdateSecretMetaReq,
	data *updateSecretData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeSecret, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *SecretUC) prepareUpdatingSecretMeta(
	req *secretdto.UpdateSecretMetaReq,
	data *updateSecretData,
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

func (uc *SecretUC) persistSecretMeta(
	ctx context.Context,
	db database.IDB,
	data *updateSecretData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
