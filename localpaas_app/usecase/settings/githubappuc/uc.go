package githubappuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GithubAppUC struct {
	*settings.BaseSettingUC
}

func NewGithubAppUC(
	baseSettingUC *settings.BaseSettingUC,
) *GithubAppUC {
	return &GithubAppUC{
		BaseSettingUC: baseSettingUC,
	}
}
