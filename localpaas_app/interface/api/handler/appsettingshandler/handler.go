package appsettingshandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/appbasehandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appsettingsuc"
)

type Handler struct {
	*appbasehandler.Handler
	appSettingsUC *appsettingsuc.UC
}

func New(
	baseHandler *appbasehandler.Handler,
	appSettingsUC *appsettingsuc.UC,
) *Handler {
	return &Handler{
		Handler:       baseHandler,
		appSettingsUC: appSettingsUC,
	}
}
