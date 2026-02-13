package awsuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type AWSUC struct {
	*settings.BaseSettingUC
}

func NewAWSUC(
	baseSettingUC *settings.BaseSettingUC,
) *AWSUC {
	return &AWSUC{
		BaseSettingUC: baseSettingUC,
	}
}
