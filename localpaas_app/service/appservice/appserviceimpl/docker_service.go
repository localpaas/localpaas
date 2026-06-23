package appserviceimpl

import (
	"context"

	"github.com/moby/moby/api/types/swarm"
	"github.com/moby/moby/client"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/services/docker"
)

func (s *service) ServiceInspect(
	ctx context.Context,
	serviceID string,
	caching bool,
) (*swarm.Service, error) {
	// TODO: handle caching flag

	resp, err := s.dockerManager.ServiceInspect(ctx, serviceID)
	if err != nil {
		return nil, apperrors.New(err)
	}

	return &resp.Service, nil
}

func (s *service) ServiceUpdate(
	ctx context.Context,
	serviceID string,
	version *swarm.Version,
	service *swarm.ServiceSpec,
	options ...docker.ServiceUpdateOption,
) (*client.ServiceUpdateResult, error) {
	resp, err := s.dockerManager.ServiceUpdate(ctx, serviceID, version, service, options...)
	if err != nil {
		return nil, apperrors.New(err)
	}
	return resp, nil
}
