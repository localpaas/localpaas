package imagebuilduc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ImageBuildUC struct {
	*settings.BaseSettingUC
}

func NewImageBuildUC(
	baseSettingUC *settings.BaseSettingUC,
) *ImageBuildUC {
	return &ImageBuildUC{
		BaseSettingUC: baseSettingUC,
	}
}
