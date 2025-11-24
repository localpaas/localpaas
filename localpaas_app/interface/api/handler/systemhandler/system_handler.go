package systemhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/nginxuc"
)

type SystemHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	nginxUC     *nginxuc.NginxUC
}

func NewSystemHandler(
	authHandler *authhandler.AuthHandler,
	nginxUC *nginxuc.NginxUC,
) *SystemHandler {
	hdl := &SystemHandler{
		authHandler: authHandler,
		nginxUC:     nginxUC,
	}
	return hdl
}
