package usersettingshandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
)

type UserSettingsHandler struct {
	*basesettinghandler.BaseSettingHandler
}

func NewUserSettingsHandler(
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
) *UserSettingsHandler {
	return &UserSettingsHandler{
		BaseSettingHandler: baseSettingHandler,
	}
}
