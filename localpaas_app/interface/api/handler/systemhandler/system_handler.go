package systemhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/lpappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/nginxuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/syserroruc"
)

type SystemHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	sysErrorUC  *syserroruc.SysErrorUC
	lpAppUC     *lpappuc.LpAppUC
	nginxUC     *nginxuc.NginxUC
}

func NewSystemHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	sysErrorUC *syserroruc.SysErrorUC,
	lpAppUC *lpappuc.LpAppUC,
	nginxUC *nginxuc.NginxUC,
) *SystemHandler {
	return &SystemHandler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		sysErrorUC:  sysErrorUC,
		lpAppUC:     lpAppUC,
		nginxUC:     nginxUC,
	}
}
