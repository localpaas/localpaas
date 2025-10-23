package apikeyhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apikeyuc"
)

type APIKeyHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	apiKeyUC    *apikeyuc.APIKeyUC
}

func NewAPIKeyHandler(
	authHandler *authhandler.AuthHandler,
	apiKeyUC *apikeyuc.APIKeyUC,
) *APIKeyHandler {
	hdl := &APIKeyHandler{
		authHandler: authHandler,
		apiKeyUC:    apiKeyUC,
	}
	return hdl
}
