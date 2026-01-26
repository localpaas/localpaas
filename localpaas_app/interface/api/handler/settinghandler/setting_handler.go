package settinghandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
)

type SettingHandler struct {
	*basesettinghandler.BaseSettingHandler
}

func NewSettingHandler(
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
) *SettingHandler {
	return &SettingHandler{
		BaseSettingHandler: baseSettingHandler,
	}
}
