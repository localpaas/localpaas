package sshkeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) ListSSHKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.ListSSHKeyReq,
) (*sshkeydto.ListSSHKeyResp, error) {
	req.Type = currentSettingType
	resp, err := providers.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &providers.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := sshkeydto.TransformSSHKeys(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sshkeydto.ListSSHKeyResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
