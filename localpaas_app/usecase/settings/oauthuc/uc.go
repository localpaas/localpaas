package oauthuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type OAuthUC struct {
	*settings.BaseSettingUC
}

func NewOAuthUC(
	baseSettingUC *settings.BaseSettingUC,
) *OAuthUC {
	return &OAuthUC{
		BaseSettingUC: baseSettingUC,
	}
}
