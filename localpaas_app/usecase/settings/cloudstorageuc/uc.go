package cloudstorageuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CloudStorageUC struct {
	*settings.BaseSettingUC
}

func NewCloudStorageUC(
	baseSettingUC *settings.BaseSettingUC,
) *CloudStorageUC {
	return &CloudStorageUC{
		BaseSettingUC: baseSettingUC,
	}
}
