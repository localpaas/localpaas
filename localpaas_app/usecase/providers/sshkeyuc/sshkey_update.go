package sshkeyuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) UpdateSSHKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.UpdateSSHKeyReq,
) (*sshkeydto.UpdateSSHKeyResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &providers.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: gofn.PtrValueOrEmpty(req.Name),
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *providers.UpdateSettingData,
			pData *providers.PersistingSettingData,
		) error {
			setting := pData.Setting
			if req.Name != nil {
				setting.Name = *req.Name
			}
			if req.PrivateKey != nil { //nolint:nestif
				sshKey, err := setting.AsSSHKey()
				if err != nil {
					return apperrors.Wrap(err)
				}
				if sshKey == nil {
					sshKey = &entity.SSHKey{}
				}
				if req.PrivateKey != nil {
					sshKey.PrivateKey = entity.NewEncryptedField(*req.PrivateKey)
				}
				if req.Passphrase != nil {
					sshKey.Passphrase = entity.NewEncryptedField(*req.Passphrase)
				}
				if err = setting.SetData(sshKey); err != nil {
					return apperrors.Wrap(err)
				}
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sshkeydto.UpdateSSHKeyResp{}, nil
}
