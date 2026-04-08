package clusterhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/imageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/networkuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/nodeuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc"
)

type Handler struct {
	*handler.BaseHandler
	authHandler *authhandler.Handler
	nodeUC      *nodeuc.UC
	volumeUC    *volumeuc.UC
	imageUC     *imageuc.UC
	networkUC   *networkuc.UC
}

func New(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.Handler,
	nodeUC *nodeuc.UC,
	volumeUC *volumeuc.UC,
	imageUC *imageuc.UC,
	networkUC *networkuc.UC,
) *Handler {
	return &Handler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		nodeUC:      nodeUC,
		volumeUC:    volumeUC,
		imageUC:     imageUC,
		networkUC:   networkUC,
	}
}
