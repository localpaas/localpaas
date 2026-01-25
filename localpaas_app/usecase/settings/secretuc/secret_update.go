package secretuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
)

func (uc *SecretUC) UpdateSecret(
	ctx context.Context,
	auth *basedto.Auth,
	req *secretdto.UpdateSecretReq,
) (*secretdto.UpdateSecretResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &settings.UpdateSettingData{
		SettingRepo: uc.settingRepo,
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *settings.UpdateSettingData,
			pData *settings.PersistingSettingData,
		) error {
			setting := pData.Setting
			secret, err := setting.AsSecret()
			if err != nil {
				return apperrors.Wrap(err)
			}
			if secret == nil {
				secret = &entity.Secret{}
			}
			secret.Value = entity.NewEncryptedField(req.Value)
			secret.Base64 = req.Base64
			if err = setting.SetData(secret); err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &secretdto.UpdateSecretResp{}, nil
}
