package providershandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
)

type ProvidersHandler struct {
	*basesettinghandler.BaseSettingHandler
	authHandler *authhandler.AuthHandler
}

func NewProvidersHandler(
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
	authHandler *authhandler.AuthHandler,
) *ProvidersHandler {
	return &ProvidersHandler{
		BaseSettingHandler: baseSettingHandler,
		authHandler:        authHandler,
	}
}
