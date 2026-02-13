package repowebhookuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type RepoWebhookUC struct {
	*settings.BaseSettingUC
}

func NewRepoWebhookUC(
	baseSettingUC *settings.BaseSettingUC,
) *RepoWebhookUC {
	return &RepoWebhookUC{
		BaseSettingUC: baseSettingUC,
	}
}
