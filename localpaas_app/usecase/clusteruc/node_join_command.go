package clusteruc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
)

func (uc *ClusterUC) GetNodeJoinCommand(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.GetNodeJoinCommandReq,
) (*clusterdto.GetNodeJoinCommandResp, error) {
	data := &joinNodeCommandData{}
	err := uc.loadGetNodeJoinCommandData(ctx, req, data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	command := fmt.Sprintf("docker swarm join --token %s %s", data.JoinToken, data.PreferManagerAddr)
	return &clusterdto.GetNodeJoinCommandResp{
		Data: &clusterdto.GetNodeJoinCommandDataResp{
			Command: command,
		},
	}, nil
}

type joinNodeCommandData struct {
	JoinToken         string
	PreferManagerAddr string
}

func (uc *ClusterUC) loadGetNodeJoinCommandData(
	ctx context.Context,
	req *clusterdto.GetNodeJoinCommandReq,
	data *joinNodeCommandData,
) error {
	// Find join token from the cluster
	theSwarm, err := uc.dockerManager.SwarmInspect(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}

	joinToken := gofn.If(req.JoinAsManager, theSwarm.JoinTokens.Manager, theSwarm.JoinTokens.Worker) //nolint
	if joinToken == "" {
		return apperrors.Wrap(apperrors.ErrDockerJoinTokenNotFound)
	}
	data.JoinToken = joinToken

	// List all manager nodes to get the addr to join new node
	managerNodes, err := uc.dockerManager.NodeList(ctx, func(opts *swarm.NodeListOptions) {
		opts.Filters = filters.NewArgs(filters.Arg("role", "manager"))
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	var leaderAddr, managerAddr string
	for _, node := range managerNodes {
		mgrStatus := node.ManagerStatus
		if mgrStatus.Reachability == swarm.ReachabilityReachable {
			managerAddr = mgrStatus.Addr
			if mgrStatus.Leader {
				leaderAddr = mgrStatus.Addr
			}
		}
	}
	data.PreferManagerAddr = gofn.Coalesce(leaderAddr, managerAddr)
	if data.PreferManagerAddr == "" {
		return apperrors.Wrap(apperrors.ErrDockerActiveManagerNodeNotFound)
	}

	return nil
}
