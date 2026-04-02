package sslcertuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type SSLCertUC struct {
	*settings.BaseSettingUC
}

func NewSSLCertUC(
	baseSettingUC *settings.BaseSettingUC,
) *SSLCertUC {
	return &SSLCertUC{
		BaseSettingUC: baseSettingUC,
	}
}
