package imserviceuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type IMServiceUC struct {
	*settings.BaseSettingUC
}

func NewIMServiceUC(
	baseSettingUC *settings.BaseSettingUC,
) *IMServiceUC {
	return &IMServiceUC{
		BaseSettingUC: baseSettingUC,
	}
}
