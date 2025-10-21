package apphandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc"
)

type AppHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	appUC       *appuc.AppUC
}

func NewAppHandler(
	authHandler *authhandler.AuthHandler,
	appUC *appuc.AppUC,
) *AppHandler {
	hdl := &AppHandler{
		authHandler: authHandler,
		appUC:       appUC,
	}
	return hdl
}
