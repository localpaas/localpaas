package apphandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc"
)

type AppHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	appUC       *appuc.AppUC
	secretUC    *secretuc.SecretUC
}

func NewAppHandler(
	authHandler *authhandler.AuthHandler,
	appUC *appuc.AppUC,
	secretUC *secretuc.SecretUC,
) *AppHandler {
	hdl := &AppHandler{
		authHandler: authHandler,
		appUC:       appUC,
		secretUC:    secretUC,
	}
	return hdl
}
