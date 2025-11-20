package registryauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) UpdateRegistryAuthMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.UpdateRegistryAuthMetaReq,
) (*registryauthdto.UpdateRegistryAuthMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		authData := &updateRegistryAuthData{}
		err := uc.loadRegistryAuthDataForUpdateMeta(ctx, db, req, authData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingRegistryAuthMeta(req, authData)
		return uc.persistRegistryAuthMeta(ctx, db, authData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &registryauthdto.UpdateRegistryAuthMetaResp{}, nil
}

func (uc *RegistryAuthUC) loadRegistryAuthDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *registryauthdto.UpdateRegistryAuthMetaReq,
	data *updateRegistryAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeRegistryAuth, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *RegistryAuthUC) prepareUpdatingRegistryAuthMeta(
	req *registryauthdto.UpdateRegistryAuthMetaReq,
	data *updateRegistryAuthData,
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

func (uc *RegistryAuthUC) persistRegistryAuthMeta(
	ctx context.Context,
	db database.IDB,
	data *updateRegistryAuthData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
