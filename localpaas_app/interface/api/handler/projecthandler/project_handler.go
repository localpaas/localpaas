package projecthandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/secretuc"
)

type ProjectHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	projectUC   *projectuc.ProjectUC
	secretUC    *secretuc.SecretUC
	basicAuthUC *basicauthuc.BasicAuthUC
}

func NewProjectHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	projectUC *projectuc.ProjectUC,
	secretUC *secretuc.SecretUC,
	basicAuthUC *basicauthuc.BasicAuthUC,
) *ProjectHandler {
	return &ProjectHandler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		projectUC:   projectUC,
		secretUC:    secretUC,
		basicAuthUC: basicAuthUC,
	}
}
