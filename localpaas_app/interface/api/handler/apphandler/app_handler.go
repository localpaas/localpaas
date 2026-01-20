package apphandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appdeploymentuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc"
)

type AppHandler struct {
	*handler.BaseHandler
	authHandler     *authhandler.AuthHandler
	appUC           *appuc.AppUC
	appDeploymentUC *appdeploymentuc.AppDeploymentUC
	secretUC        *secretuc.SecretUC
}

func NewAppHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	appUC *appuc.AppUC,
	appDeploymentUC *appdeploymentuc.AppDeploymentUC,
	secretUC *secretuc.SecretUC,
) *AppHandler {
	return &AppHandler{
		BaseHandler:     baseHandler,
		authHandler:     authHandler,
		appUC:           appUC,
		appDeploymentUC: appDeploymentUC,
		secretUC:        secretUC,
	}
}
