package secretuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type SecretUC struct {
	*settings.BaseSettingUC
}

func NewSecretUC(
	baseSettingUC *settings.BaseSettingUC,
) *SecretUC {
	return &SecretUC{
		BaseSettingUC: baseSettingUC,
	}
}
