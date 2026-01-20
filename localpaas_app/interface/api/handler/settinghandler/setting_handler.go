package settinghandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
)

type SettingHandler struct {
	*basesettinghandler.BaseSettingHandler
	authHandler *authhandler.AuthHandler
}

func NewSettingHandler(
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
	authHandler *authhandler.AuthHandler,
) *SettingHandler {
	return &SettingHandler{
		BaseSettingHandler: baseSettingHandler,
		authHandler:        authHandler,
	}
}
