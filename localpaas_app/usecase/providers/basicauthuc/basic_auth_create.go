package basicauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
)

const (
	currentSettingType = base.SettingTypeBasicAuth
)

func (uc *BasicAuthUC) CreateBasicAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.CreateBasicAuthReq,
) (*basicauthdto.CreateBasicAuthResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		Version:       entity.CurrentBasicAuthVersion,
		PrepareCreation: func(ctx context.Context, db database.Tx, data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData) error {
			err := pData.Setting.SetData(&entity.BasicAuth{
				Username: req.Username,
				Password: entity.NewEncryptedField(req.Password),
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

	return &basicauthdto.CreateBasicAuthResp{
		Data: resp.Data,
	}, nil
}
