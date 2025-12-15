package systemhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/lpappuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/nginxuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/syserroruc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/taskuc"
)

type SystemHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	sysErrorUC  *syserroruc.SysErrorUC
	taskUC      *taskuc.TaskUC
	lpAppUC     *lpappuc.LpAppUC
	nginxUC     *nginxuc.NginxUC
}

func NewSystemHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	sysErrorUC *syserroruc.SysErrorUC,
	taskUC *taskuc.TaskUC,
	lpAppUC *lpappuc.LpAppUC,
	nginxUC *nginxuc.NginxUC,
) *SystemHandler {
	return &SystemHandler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		sysErrorUC:  sysErrorUC,
		taskUC:      taskUC,
		lpAppUC:     lpAppUC,
		nginxUC:     nginxUC,
	}
}
