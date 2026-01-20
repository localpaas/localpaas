package registryauthuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) UpdateRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.UpdateRegistryAuthReq,
) (*registryauthdto.UpdateRegistryAuthResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &providers.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *providers.UpdateSettingData,
			pData *providers.PersistingSettingData,
		) error {
			pData.Setting.Name = gofn.Coalesce(req.Name, pData.Setting.Name)
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

	return &registryauthdto.UpdateRegistryAuthResp{}, nil
}
