package cloudprovideruc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CloudProviderUC struct {
	*settings.BaseSettingUC
}

func NewCloudProviderUC(
	baseSettingUC *settings.BaseSettingUC,
) *CloudProviderUC {
	return &CloudProviderUC{
		BaseSettingUC: baseSettingUC,
	}
}
