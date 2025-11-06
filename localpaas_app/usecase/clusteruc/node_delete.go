package clusteruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
)

func (uc *ClusterUC) DeleteNode(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.DeleteNodeReq,
) (*clusterdto.DeleteNodeResp, error) {
	err := uc.dockerManager.NodeRemove(ctx, req.NodeID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &clusterdto.DeleteNodeResp{}, nil
}
