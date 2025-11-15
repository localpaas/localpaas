package clusterhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/nodeuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc"
)

type ClusterHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	nodeUC      *nodeuc.NodeUC
	volumeUC    *volumeuc.VolumeUC
}

func NewClusterHandler(
	authHandler *authhandler.AuthHandler,
	nodeUC *nodeuc.NodeUC,
	volumeUC *volumeuc.VolumeUC,
) *ClusterHandler {
	hdl := &ClusterHandler{
		authHandler: authHandler,
		nodeUC:      nodeUC,
		volumeUC:    volumeUC,
	}
	return hdl
}
