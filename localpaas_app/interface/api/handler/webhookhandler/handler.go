package webhookhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc"
)

type Handler struct {
	*handler.BaseHandler
	webhookUC *webhookuc.UC
}

func New(
	baseHandler *handler.BaseHandler,
	webhookUC *webhookuc.UC,
) *Handler {
	return &Handler{
		BaseHandler: baseHandler,
		webhookUC:   webhookUC,
	}
}
