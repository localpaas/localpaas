package secretuc

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
)

func (uc *SecretUC) CreateSecret(
	ctx context.Context,
	auth *basedto.Auth,
	req *secretdto.CreateSecretReq,
) (*secretdto.CreateSecretResp, error) {
	secretData := &createSecretData{}
	err := uc.loadSecretData(ctx, uc.db, req, secretData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingSecretData{}
	uc.preparePersistingSecret(req, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &secretdto.CreateSecretResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createSecretData struct {
}

func (uc *SecretUC) loadSecretData(
	ctx context.Context,
	db database.IDB,
	req *secretdto.CreateSecretReq,
	_ *createSecretData,
) error {
	var options []bunex.SelectQueryOption
	if req.ObjectID != "" {
		options = append(options, bunex.SelectWhere("setting.object_id = ?", req.ObjectID))
	} else {
		options = append(options, bunex.SelectWhere("setting.object_id IS NULL"))
	}

	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeSecret, req.Key, false, options...)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("Secret").
			WithMsgLog("secret '%s' already exists", req.Key)
	}

	return nil
}

type persistingSecretData struct {
	settingservice.PersistingSettingData
}

func (uc *SecretUC) preparePersistingSecret(
	req *secretdto.CreateSecretReq,
	persistingData *persistingSecretData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeSecret,
		Status:    base.SettingStatusActive,
		Name:      req.Key,
		ObjectID:  req.ObjectID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	secret := &entity.Secret{
		Key:    req.Key,
		Value:  req.Value,
		Base64: req.Base64,
	}
	setting.MustSetData(secret.MustEncrypt())

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *SecretUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingSecretData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
