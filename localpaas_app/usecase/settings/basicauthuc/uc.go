package basicauthuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type BasicAuthUC struct {
	*settings.BaseSettingUC
}

func NewBasicAuthUC(
	baseSettingUC *settings.BaseSettingUC,
) *BasicAuthUC {
	return &BasicAuthUC{
		BaseSettingUC: baseSettingUC,
	}
}
