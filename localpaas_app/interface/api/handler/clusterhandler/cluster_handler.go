package clusterhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/imageuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/nodeuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc"
)

type ClusterHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	nodeUC      *nodeuc.NodeUC
	volumeUC    *volumeuc.VolumeUC
	imageUC     *imageuc.ImageUC
}

func NewClusterHandler(
	authHandler *authhandler.AuthHandler,
	nodeUC *nodeuc.NodeUC,
	volumeUC *volumeuc.VolumeUC,
	imageUC *imageuc.ImageUC,
) *ClusterHandler {
	hdl := &ClusterHandler{
		authHandler: authHandler,
		nodeUC:      nodeUC,
		volumeUC:    volumeUC,
		imageUC:     imageUC,
	}
	return hdl
}
