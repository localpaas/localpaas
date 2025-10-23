package sshkeyhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sshkeyuc"
)

type SSHKeyHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	sshKeyUC    *sshkeyuc.SSHKeyUC
}

func NewSSHKeyHandler(
	authHandler *authhandler.AuthHandler,
	sshKeyUC *sshkeyuc.SSHKeyUC,
) *SSHKeyHandler {
	hdl := &SSHKeyHandler{
		authHandler: authHandler,
		sshKeyUC:    sshKeyUC,
	}
	return hdl
}
