package githubappuc

import (
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GithubAppUC struct {
	*settings.BaseSettingUC
	cacheAppManifestRepo cacherepository.GithubAppManifestRepo
}

func NewGithubAppUC(
	baseSettingUC *settings.BaseSettingUC,
	cacheAppManifestRepo cacherepository.GithubAppManifestRepo,
) *GithubAppUC {
	return &GithubAppUC{
		BaseSettingUC:        baseSettingUC,
		cacheAppManifestRepo: cacheAppManifestRepo,
	}
}
