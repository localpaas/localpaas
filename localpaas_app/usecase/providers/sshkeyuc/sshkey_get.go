package sshkeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) GetSSHKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.GetSSHKeyReq,
) (*sshkeydto.GetSSHKeyResp, error) {
	req.Type = currentSettingType
	setting, err := providers.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &providers.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsSSHKey().MustDecrypt()
	resp, err := sshkeydto.TransformSSHKey(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sshkeydto.GetSSHKeyResp{
		Data: resp,
	}, nil
}
