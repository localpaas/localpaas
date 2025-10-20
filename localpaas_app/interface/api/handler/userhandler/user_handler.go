package userhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc"
)

type UserHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	userUC      *useruc.UserUC
}

func NewUserHandler(authHandler *authhandler.AuthHandler, userUC *useruc.UserUC) *UserHandler {
	hdl := &UserHandler{
		authHandler: authHandler,
		userUC:      userUC,
	}
	return hdl
}
