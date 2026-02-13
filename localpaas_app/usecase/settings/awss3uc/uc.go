package awss3uc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type AWSS3UC struct {
	*settings.BaseSettingUC
}

func NewAWSS3UC(
	baseSettingUC *settings.BaseSettingUC,
) *AWSS3UC {
	return &AWSS3UC{
		BaseSettingUC: baseSettingUC,
	}
}
