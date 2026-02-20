package registryauthuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/services/docker"
)

type RegistryAuthUC struct {
	*settings.BaseSettingUC
	dockerManager docker.Manager
}

func NewRegistryAuthUC(
	baseSettingUC *settings.BaseSettingUC,
	dockerManager docker.Manager,
) *RegistryAuthUC {
	return &RegistryAuthUC{
		BaseSettingUC: baseSettingUC,
		dockerManager: dockerManager,
	}
}
