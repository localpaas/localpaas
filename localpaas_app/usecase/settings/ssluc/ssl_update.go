package ssluc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

func (uc *SslUC) UpdateSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.UpdateSslReq,
) (*ssldto.UpdateSslResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &settings.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *settings.UpdateSettingData,
			pData *settings.PersistingSettingData,
		) error {
			pData.Setting.Name = gofn.Coalesce(req.Name, pData.Setting.Name)
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

	return &ssldto.UpdateSslResp{}, nil
}
