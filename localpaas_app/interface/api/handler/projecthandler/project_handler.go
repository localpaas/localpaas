package projecthandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
)

type ProjectHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	projectUC   *projectuc.ProjectUC
}

func NewProjectHandler(
	authHandler *authhandler.AuthHandler,
	projectUC *projectuc.ProjectUC,
) *ProjectHandler {
	hdl := &ProjectHandler{
		authHandler: authHandler,
		projectUC:   projectUC,
	}
	return hdl
}
