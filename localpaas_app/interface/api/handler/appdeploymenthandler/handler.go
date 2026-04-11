package appdeploymenthandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/appbasehandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appdeploymentuc"
)

type Handler struct {
	*appbasehandler.Handler
	appDeploymentUC *appdeploymentuc.UC
}

func New(
	baseHandler *appbasehandler.Handler,
	appDeploymentUC *appdeploymentuc.UC,
) *Handler {
	return &Handler{
		Handler:         baseHandler,
		appDeploymentUC: appDeploymentUC,
	}
}
