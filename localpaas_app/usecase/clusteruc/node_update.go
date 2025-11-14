package clusteruc

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
)

func (uc *ClusterUC) UpdateNode(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.UpdateNodeReq,
) (*clusterdto.UpdateNodeResp, error) {
	node, _, err := uc.dockerManager.NodeInspect(ctx, req.NodeID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	spec := node.Spec

	if req.Name != "" {
		spec.Annotations.Name = req.Name //nolint
	}
	if req.Labels != nil {
		spec.Annotations.Labels = req.Labels //nolint
	}
	if req.Role != "" {
		spec.Role = swarm.NodeRole(req.Role)
	}
	if req.Availability != "" {
		spec.Availability = swarm.NodeAvailability(req.Availability)
	}

	err = uc.dockerManager.NodeUpdate(ctx, req.NodeID, &node.Version, &spec)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrDockerFailedUpdateNode).WithNTParam("Error", err.Error())
	}

	return &clusterdto.UpdateNodeResp{}, nil
}
