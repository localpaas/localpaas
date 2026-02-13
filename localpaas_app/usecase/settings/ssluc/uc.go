package ssluc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type SSLUC struct {
	*settings.BaseSettingUC
}

func NewSSLUC(
	baseSettingUC *settings.BaseSettingUC,
) *SSLUC {
	return &SSLUC{
		BaseSettingUC: baseSettingUC,
	}
}
