package nodeuc

import (
	"context"

	"github.com/moby/moby/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/dockerhelper"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/nodeuc/nodedto"
)

func (uc *UC) UpdateNode(
	ctx context.Context,
	auth *basedto.Auth,
	req *nodedto.UpdateNodeReq,
) (*nodedto.UpdateNodeResp, error) {
	inspect, err := uc.dockerManager.NodeInspect(ctx, req.NodeID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	node := &inspect.Node

	err = uc.verifyNodeUpdateChange(ctx, req, node)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	spec := &node.Spec

	if req.Name != "" {
		spec.Annotations.Name = req.Name //nolint
	}
	spec.Labels = dockerhelper.ApplyUserLabels(spec.Labels, req.Labels)
	if req.Role != "" {
		spec.Role = swarm.NodeRole(req.Role)
	}
	if req.Availability != "" {
		spec.Availability = swarm.NodeAvailability(req.Availability)
	}

	_, err = uc.dockerManager.NodeUpdate(ctx, req.NodeID, &node.Version, spec)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &nodedto.UpdateNodeResp{}, nil
}

func (uc *UC) verifyNodeUpdateChange(
	ctx context.Context,
	req *nodedto.UpdateNodeReq,
	node *swarm.Node,
) error {
	if uint64(req.UpdateVer) != node.Version.Index { //nolint:gosec
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}

	spec := &node.Spec

	roleDemoting := swarm.NodeRole(req.Role) == swarm.NodeRoleWorker && spec.Role == swarm.NodeRoleManager
	availabilityLosing := swarm.NodeAvailability(req.Availability) != swarm.NodeAvailabilityActive &&
		spec.Availability == swarm.NodeAvailabilityActive

	if roleDemoting || availabilityLosing {
		tasks, err := uc.lpAppService.GetLpAppTasks(ctx)
		if err != nil {
			return apperrors.Wrap(err)
		}
		allNodes := make(map[string]*swarm.Task)
		for i := range tasks {
			allNodes[tasks[i].NodeID] = &tasks[i]
		}
		if len(allNodes) == 1 && allNodes[node.ID] != nil {
			return apperrors.New(apperrors.ErrNodeRequiredByLocalPaaSApp).WithDisplayLevelHigh()
		}
	}
	return nil
}
