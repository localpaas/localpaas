package usersettingshandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc"
)

type UserSettingsHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	apiKeyUC    *apikeyuc.APIKeyUC
}

func NewUserSettingsHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	apiKeyUC *apikeyuc.APIKeyUC,
) *UserSettingsHandler {
	return &UserSettingsHandler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		apiKeyUC:    apiKeyUC,
	}
}
