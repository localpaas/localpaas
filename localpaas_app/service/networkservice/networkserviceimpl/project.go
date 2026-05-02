package networkserviceimpl

import (
	"context"

	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/services/docker"
)

func (s *service) GetProjectNetwork(
	ctx context.Context,
	project *entity.Project,
) (*network.Inspect, error) {
	// Create a default network for the project apps
	inspect, err := s.dockerManager.NetworkInspect(ctx, project.GetDefaultNetworkName())
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return &inspect.Network, nil
}

func (s *service) CreateProjectNetwork(
	ctx context.Context,
	project *entity.Project,
) (*client.NetworkCreateResult, error) {
	// Create a default network for the project apps
	resp, err := s.dockerManager.NetworkCreate(ctx, project.GetDefaultNetworkName(),
		func(opts *client.NetworkCreateOptions) {
			opts.Driver = docker.NetworkDriverOverlay
			opts.Scope = docker.NetworkScopeSwarm
			opts.Attachable = true
			opts.Labels = map[string]string{
				docker.StackLabelNamespace: project.Key,
			}
		})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func (s *service) RemoveProjectNetwork(
	ctx context.Context,
	project *entity.Project,
) error {
	_, err := s.dockerManager.NetworkRemove(ctx, project.GetDefaultNetworkName())
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
