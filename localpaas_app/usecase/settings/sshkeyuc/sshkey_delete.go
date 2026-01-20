package sshkeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) DeleteSSHKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.DeleteSSHKeyReq,
) (*sshkeydto.DeleteSSHKeyResp, error) {
	req.Type = currentSettingType
	_, err := settings.DeleteSetting(ctx, uc.db, &req.DeleteSettingReq, &settings.DeleteSettingData{
		SettingRepo:              uc.settingRepo,
		ProjectSharedSettingRepo: uc.projectSharedSettingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sshkeydto.DeleteSSHKeyResp{}, nil
}
