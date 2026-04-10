package traefikserviceimpl

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

func (s *service) GetTraefikSwarmService(ctx context.Context) (*swarm.Service, error) {
	service, err := s.dockerManager.ServiceGetByName(ctx, base.LocalpaasTraefikServiceName)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return service, nil
}

func (s *service) RestartTraefikSwarmService(ctx context.Context) error {
	service, err := s.GetTraefikSwarmService(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.dockerManager.ServiceForceUpdate(ctx, service.ID)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
