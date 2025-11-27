package registryauthuc

import (
	"context"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) UpdateRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.UpdateRegistryAuthReq,
) (*registryauthdto.UpdateRegistryAuthResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		registryAuthData := &updateRegistryAuthData{}
		err := uc.loadRegistryAuthDataForUpdate(ctx, db, req, registryAuthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingRegistryAuthData{}
		uc.prepareUpdatingRegistryAuth(req.RegistryAuthBaseReq, registryAuthData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &registryauthdto.UpdateRegistryAuthResp{}, nil
}

type updateRegistryAuthData struct {
	Setting *entity.Setting
}

func (uc *RegistryAuthUC) loadRegistryAuthDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *registryauthdto.UpdateRegistryAuthReq,
	data *updateRegistryAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeRegistryAuth, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if req.Name != "" && !strings.EqualFold(setting.Name, req.Name) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeRegistryAuth, req.Name, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("RegistryAuth").
				WithMsgLog("registry auth '%s' already exists", conflictSetting.Name)
		}
	}

	return nil
}

func (uc *RegistryAuthUC) prepareUpdatingRegistryAuth(
	req *registryauthdto.RegistryAuthBaseReq,
	data *updateRegistryAuthData,
	persistingData *persistingRegistryAuthData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	if req.Name != "" {
		setting.Name = req.Name
	}
	setting.Kind = req.Address

	registryAuth := &entity.RegistryAuth{
		Username: req.Username,
		Password: entity.NewEncryptedField(req.Password),
		Address:  req.Address,
	}
	setting.MustSetData(registryAuth)

	setting.UpdatedAt = timeNow
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
