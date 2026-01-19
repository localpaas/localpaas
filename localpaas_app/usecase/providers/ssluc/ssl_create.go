package ssluc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc/ssldto"
)

const (
	currentSettingType    = base.SettingTypeSSL
	currentSettingVersion = entity.CurrentSslVersion
)

func (uc *SslUC) CreateSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.CreateSslReq,
) (*ssldto.CreateSslResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		Version:       currentSettingVersion,
		PrepareCreation: func(ctx context.Context, db database.Tx, data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData) error {
			pData.Setting.Kind = string(req.Provider)
			err := pData.Setting.SetData(&entity.Ssl{
				Certificate: req.Certificate,
				PrivateKey:  entity.NewEncryptedField(req.PrivateKey),
				KeySize:     req.KeySize,
				Provider:    req.Provider,
				Email:       req.Email,
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

	return &ssldto.CreateSslResp{
		Data: resp.Data,
	}, nil
}
