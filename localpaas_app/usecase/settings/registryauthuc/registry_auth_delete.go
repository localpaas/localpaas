package registryauthuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) DeleteRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.DeleteRegistryAuthReq,
) (*registryauthdto.DeleteRegistryAuthResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		registryAuthData := &deleteRegistryAuthData{}
		err := uc.loadRegistryAuthDataForDelete(ctx, db, req, registryAuthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingRegistryAuthData{}
		uc.prepareDeletingRegistryAuth(registryAuthData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &registryauthdto.DeleteRegistryAuthResp{}, nil
}

type deleteRegistryAuthData struct {
	Setting *entity.Setting
}

func (uc *RegistryAuthUC) loadRegistryAuthDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *registryauthdto.DeleteRegistryAuthReq,
	data *deleteRegistryAuthData,
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

func (uc *RegistryAuthUC) prepareDeletingRegistryAuth(
	data *deleteRegistryAuthData,
	persistingData *persistingRegistryAuthData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
