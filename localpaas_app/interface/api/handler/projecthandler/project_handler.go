package projecthandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/secretuc"
)

type ProjectHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	projectUC   *projectuc.ProjectUC
	secretUC    *secretuc.SecretUC
}

func NewProjectHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	projectUC *projectuc.ProjectUC,
	secretUC *secretuc.SecretUC,
) *ProjectHandler {
	return &ProjectHandler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		projectUC:   projectUC,
		secretUC:    secretUC,
	}
}
