package imagehandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/binobjectuc"
)

type Handler struct {
	*handler.BaseHandler
	authHandler *authhandler.Handler
	binObjectUC *binobjectuc.UC
}

func New(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.Handler,
	binObjectUC *binobjectuc.UC,
) *Handler {
	return &Handler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		binObjectUC: binObjectUC,
	}
}
