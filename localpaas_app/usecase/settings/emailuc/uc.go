package emailuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type EmailUC struct {
	*settings.BaseSettingUC
}

func NewEmailUC(
	baseSettingUC *settings.BaseSettingUC,
) *EmailUC {
	return &EmailUC{
		BaseSettingUC: baseSettingUC,
	}
}
