package accesstokenuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type AccessTokenUC struct {
	*settings.BaseSettingUC
}

func NewAccessTokenUC(
	baseSettingUC *settings.BaseSettingUC,
) *AccessTokenUC {
	return &AccessTokenUC{
		BaseSettingUC: baseSettingUC,
	}
}
