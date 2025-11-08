package settingshandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc"
)

type SettingsHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	oauthUC     *oauthuc.OAuthUC
}

func NewSettingsHandler(
	authHandler *authhandler.AuthHandler,
	oauthUC *oauthuc.OAuthUC,
) *SettingsHandler {
	hdl := &SettingsHandler{
		authHandler: authHandler,
		oauthUC:     oauthUC,
	}
	return hdl
}
