package clusterhandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc"
)

type ClusterHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	clusterUC   *clusteruc.ClusterUC
}

func NewClusterHandler(
	authHandler *authhandler.AuthHandler,
	clusterUC *clusteruc.ClusterUC,
) *ClusterHandler {
	hdl := &ClusterHandler{
		authHandler: authHandler,
		clusterUC:   clusterUC,
	}
	return hdl
}
