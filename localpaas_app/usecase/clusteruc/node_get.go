package clusteruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
)

func (uc *ClusterUC) GetNode(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.GetNodeReq,
) (*clusterdto.GetNodeResp, error) {
	node, _, err := uc.dockerManager.NodeInspect(ctx, req.NodeID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &clusterdto.GetNodeResp{
		Data: clusterdto.TransformNode(node, true),
	}, nil
}
