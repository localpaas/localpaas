package networkservice

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	networkLabelProjectKey = "localpaas.project.key"
)

func (s *networkService) CreateProjectNetwork(ctx context.Context, project *entity.Project) (
	*network.CreateResponse, error) {
	// Create a default network for the project apps
	net, err := s.dockerManager.NetworkCreate(ctx, project.GetDefaultNetworkName(), func(opts *network.CreateOptions) {
		opts.Driver = swarmNetworkDriver
		opts.Scope = swarmNetworkScope
		opts.Attachable = true
		opts.Labels = map[string]string{
			networkLabelProjectKey:     project.Key,
			docker.StackLabelNamespace: project.Key,
		}
	})
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return net, nil
}

func (s *networkService) ListProjectNetworks(ctx context.Context, project *entity.Project) (
	[]network.Summary, error) {
	res, err := s.dockerManager.NetworkList(ctx, func(opts *network.ListOptions) {
		opts.Filters = filters.NewArgs(
			filters.Arg("label", networkLabelProjectKey+"="+project.Key),
		)
	})
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return res, nil
}

func (s *networkService) RemoveProjectNetwork(ctx context.Context, project *entity.Project) error {
	err := s.dockerManager.NetworkRemove(ctx, project.GetDefaultNetworkName())
	if err != nil {
		return apperrors.NewInfra(err)
	}
	return nil
}
