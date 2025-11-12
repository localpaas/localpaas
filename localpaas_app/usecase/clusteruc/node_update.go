package clusteruc

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
	"github.com/localpaas/localpaas/services/docker"
)

func (uc *ClusterUC) UpdateNode(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.UpdateNodeReq,
) (*clusterdto.UpdateNodeResp, error) {
	err := uc.dockerManager.NodeUpdate(ctx, req.NodeID, docker.VersionAuto, &swarm.NodeSpec{
		Annotations: swarm.Annotations{
			Name:   req.Name,
			Labels: req.Labels,
		},
		Role:         swarm.NodeRole(req.Role),
		Availability: swarm.NodeAvailability(req.Availability),
	})
	if err != nil {
		return nil, apperrors.New(apperrors.ErrDockerFailedUpdateNode).WithNTParam("Error", err.Error())
	}

	return &clusterdto.UpdateNodeResp{}, nil
}
