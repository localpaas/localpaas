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
	authHandler *authhandler.AuthHandler,
	apiKeyUC *apikeyuc.APIKeyUC,
) *UserSettingsHandler {
	hdl := &UserSettingsHandler{
		authHandler: authHandler,
		apiKeyUC:    apiKeyUC,
	}
	return hdl
}
