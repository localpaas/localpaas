package systemsettingshandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
)

type SystemSettingsHandler struct {
	*basesettinghandler.BaseSettingHandler
}

func NewSystemSettingsHandler(
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
) *SystemSettingsHandler {
	return &SystemSettingsHandler{
		BaseSettingHandler: baseSettingHandler,
	}
}
