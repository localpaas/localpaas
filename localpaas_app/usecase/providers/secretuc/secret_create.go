package secretuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/secretuc/secretdto"
)

const (
	currentSettingType    = base.SettingTypeSecret
	currentSettingVersion = entity.CurrentSecretVersion
)

func (uc *SecretUC) CreateSecret(
	ctx context.Context,
	auth *basedto.Auth,
	req *secretdto.CreateSecretReq,
) (*secretdto.CreateSecretResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Key,
		Version:       currentSettingVersion,
		PrepareCreation: func(ctx context.Context, db database.Tx, data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData) error {
			err := pData.Setting.SetData(&entity.Secret{
				Key:    req.Key,
				Value:  entity.NewEncryptedField(req.Value),
				Base64: req.Base64,
			})
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &secretdto.CreateSecretResp{
		Data: resp.Data,
	}, nil
}
