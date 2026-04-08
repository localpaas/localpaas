package apphandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appdeploymentuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc"
)

type Handler struct {
	*basesettinghandler.Handler
	appUC           *appuc.UC
	appDeploymentUC *appdeploymentuc.UC
	secretUC        *secretuc.UC
}

func New(
	baseSettingHandler *basesettinghandler.Handler,
	appUC *appuc.UC,
	appDeploymentUC *appdeploymentuc.UC,
	secretUC *secretuc.UC,
) *Handler {
	return &Handler{
		Handler:         baseSettingHandler,
		appUC:           appUC,
		appDeploymentUC: appDeploymentUC,
		secretUC:        secretUC,
	}
}
