package apikeyuc

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

const (
	keyLen           = 16
	secretLen        = 32
	saltLen          = 8
	hashingKeyLen    = 32
	hashingIteration = 1
)

func (uc *APIKeyUC) CreateAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.CreateAPIKeyReq,
) (*apikeydto.CreateAPIKeyResp, error) {
	apiKeyData := &createAPIKeyData{}
	err := uc.loadAPIKeyData(ctx, uc.db, auth, req, apiKeyData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingAPIKeyData{}
	err = uc.preparePersistingAPIKey(req, apiKeyData, persistingData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.CreateAPIKeyResp{
		Data: &apikeydto.APIKeyDataResp{
			ID:        persistingData.UpsertingSettings[0].ID,
			KeyID:     apiKeyData.KeyID,
			SecretKey: apiKeyData.SecretKey,
		},
	}, nil
}

type createAPIKeyData struct {
	ActingUser *entity.User
	KeyID      string
	SecretKey  string
}

func (uc *APIKeyUC) loadAPIKeyData(
	ctx context.Context,
	db database.IDB,
	auth *basedto.Auth,
	_ *apikeydto.CreateAPIKeyReq,
	data *createAPIKeyData,
) error {
	data.ActingUser = auth.User.User

	// Generate key and secret
	keyID, secretKey := gofn.RandTokenAsHex(keyLen), gofn.RandTokenAsHex(secretLen)
	data.KeyID = keyID
	data.SecretKey = secretKey

	// Make sure there is no duplicated key in the db
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeAPIKey, keyID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("APIKey").
			WithMsgLog("API key '%s' setting already exists", keyID)
	}

	return nil
}

type persistingAPIKeyData struct {
	settingservice.PersistingSettingData
}

func (uc *APIKeyUC) preparePersistingAPIKey(
	req *apikeydto.CreateAPIKeyReq,
	data *createAPIKeyData,
	persistingData *persistingAPIKeyData,
) error {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeAPIKey,
		Status:    base.SettingStatusActive,
		ObjectID:  data.ActingUser.ID,
		Name:      data.KeyID,
		ExpireAt:  req.Expiration,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	apiKey := &entity.APIKey{
		SecretKey:    data.SecretKey,
		AccessAction: req.AccessAction,
	}
	err := apiKey.Hash()
	if err != nil {
		return apperrors.Wrap(err)
	}
	setting.MustSetData(apiKey)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	return nil
}

func (uc *APIKeyUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingAPIKeyData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
