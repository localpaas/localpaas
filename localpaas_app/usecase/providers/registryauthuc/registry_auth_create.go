package registryauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc/registryauthdto"
)

const (
	currentSettingType    = base.SettingTypeRegistryAuth
	currentSettingVersion = entity.CurrentRegistryAuthVersion
)

func (uc *RegistryAuthUC) CreateRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.CreateRegistryAuthReq,
) (*registryauthdto.CreateRegistryAuthResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		Version:       currentSettingVersion,
		PrepareCreation: func(
			ctx context.Context,
			db database.Tx,
			data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData,
		) error {
			pData.Setting.Kind = req.Address
			err := pData.Setting.SetData(&entity.RegistryAuth{
				Username: req.Username,
				Password: entity.NewEncryptedField(req.Password),
				Address:  req.Address,
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

	return &registryauthdto.CreateRegistryAuthResp{
		Data: resp.Data,
	}, nil
}
