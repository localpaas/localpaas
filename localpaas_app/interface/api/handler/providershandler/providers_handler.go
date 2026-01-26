package providershandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
)

type ProvidersHandler struct {
	*basesettinghandler.BaseSettingHandler
}

func NewProvidersHandler(
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
) *ProvidersHandler {
	return &ProvidersHandler{
		BaseSettingHandler: baseSettingHandler,
	}
}
