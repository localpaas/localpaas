package apikeyuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type APIKeyUC struct {
	*settings.BaseSettingUC
}

func NewAPIKeyUC(
	baseSettingUC *settings.BaseSettingUC,
) *APIKeyUC {
	return &APIKeyUC{
		BaseSettingUC: baseSettingUC,
	}
}
