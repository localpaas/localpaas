package sshkeyuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type SSHKeyUC struct {
	*settings.BaseSettingUC
}

func NewSSHKeyUC(
	baseSettingUC *settings.BaseSettingUC,
) *SSHKeyUC {
	return &SSHKeyUC{
		BaseSettingUC: baseSettingUC,
	}
}
