package usersettingshandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
)

type UserSettingsHandler struct {
	*basesettinghandler.BaseSettingHandler
	authHandler *authhandler.AuthHandler
}

func NewUserSettingsHandler(
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
	authHandler *authhandler.AuthHandler,
) *UserSettingsHandler {
	return &UserSettingsHandler{
		BaseSettingHandler: baseSettingHandler,
		authHandler:        authHandler,
	}
}
