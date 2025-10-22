package projectenvhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc"
)

type ProjectEnvHandler struct {
	*handler.BaseHandler
	authHandler  *authhandler.AuthHandler
	projectEnvUC *projectenvuc.ProjectEnvUC
}

func NewProjectEnvHandler(
	authHandler *authhandler.AuthHandler,
	projectEnvUC *projectenvuc.ProjectEnvUC,
) *ProjectEnvHandler {
	hdl := &ProjectEnvHandler{
		authHandler:  authHandler,
		projectEnvUC: projectEnvUC,
	}
	return hdl
}
