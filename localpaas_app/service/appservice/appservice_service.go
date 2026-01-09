package appservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/services/docker"
)

func (s *appService) ServiceInspect(
	ctx context.Context,
	serviceID string,
) (*swarm.Service, error) {
	service, err := s.dockerManager.ServiceInspect(ctx, serviceID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return service, nil
}

func (s *appService) ServiceUpdate(
	ctx context.Context,
	serviceID string,
	version *swarm.Version,
	service *swarm.ServiceSpec,
	options ...docker.ServiceUpdateOption,
) (*swarm.ServiceUpdateResponse, error) {
	resp, err := s.dockerManager.ServiceUpdate(ctx, serviceID, version, service, options...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
