package apppreviewhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/appbasehandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apppreviewuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc"
)

type Handler struct {
	*appbasehandler.Handler
	appUC        *appuc.UC
	appPreviewUC *apppreviewuc.UC
}

func New(
	baseHandler *appbasehandler.Handler,
	appUC *appuc.UC,
	appPreviewUC *apppreviewuc.UC,
) *Handler {
	return &Handler{
		Handler:      baseHandler,
		appUC:        appUC,
		appPreviewUC: appPreviewUC,
	}
}
