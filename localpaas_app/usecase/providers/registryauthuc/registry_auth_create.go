package registryauthuc

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) CreateRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.CreateRegistryAuthReq,
) (*registryauthdto.CreateRegistryAuthResp, error) {
	registryAuthData := &createRegistryAuthData{}
	err := uc.loadRegistryAuthData(ctx, uc.db, req, registryAuthData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingRegistryAuthData{}
	uc.preparePersistingRegistryAuth(req.RegistryAuthBaseReq, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &registryauthdto.CreateRegistryAuthResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createRegistryAuthData struct {
}

func (uc *RegistryAuthUC) loadRegistryAuthData(
	ctx context.Context,
	db database.IDB,
	req *registryauthdto.CreateRegistryAuthReq,
	_ *createRegistryAuthData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeRegistryAuth, req.Name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("RegistryAuth").
			WithMsgLog("registry auth setting '%s' already exists", req.Name)
	}

	return nil
}

type persistingRegistryAuthData struct {
	settingservice.PersistingSettingData
}

func (uc *RegistryAuthUC) preparePersistingRegistryAuth(
	req *registryauthdto.RegistryAuthBaseReq,
	persistingData *persistingRegistryAuthData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeRegistryAuth,
		Status:    base.SettingStatusActive,
		Kind:      req.Address,
		Name:      req.Name,
		Version:   entity.CurrentRegistryAuthVersion,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	registryAuth := &entity.RegistryAuth{
		Username: req.Username,
		Password: entity.NewEncryptedField(req.Password),
		Address:  req.Address,
	}
	setting.MustSetData(registryAuth)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *RegistryAuthUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingRegistryAuthData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
