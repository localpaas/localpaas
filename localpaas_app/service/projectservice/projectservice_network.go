package projectservice

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

const (
	networkDriver = "overlay"
	networkScope  = "swarm"

	networkLabelProjectSlug = "localpaas.project.slug"
)

func (s *projectService) CreateProjectNetworks(ctx context.Context, project *entity.Project) (
	*network.CreateResponse, error) {
	// Create a default network for the project apps
	net, err := s.dockerManager.NetworkCreate(ctx, project.GetDefaultNetworkName(), func(opts *network.CreateOptions) {
		opts.Driver = networkDriver
		opts.Scope = networkScope
		opts.Attachable = true
		opts.Labels = map[string]string{
			networkLabelProjectSlug: project.Slug,
		}
	})
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return net, nil
}

func (s *projectService) ListProjectNetworks(ctx context.Context, project *entity.Project) (
	[]network.Summary, error) {
	res, err := s.dockerManager.NetworkList(ctx, func(opts *network.ListOptions) {
		opts.Filters = filters.NewArgs(
			filters.Arg("label", networkLabelProjectSlug+"="+project.Slug),
		)
	})
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return res, nil
}
