package lpappservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	lpAppServiceName = "localpaas_app"
)

func (s *lpAppService) GetLpAppSwarmService(ctx context.Context) (*swarm.Service, error) {
	service, err := s.dockerManager.ServiceGetByName(ctx, lpAppServiceName)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return service, nil
}

func (s *lpAppService) RestartLpAppSwarmService(ctx context.Context) error {
	service, err := s.GetLpAppSwarmService(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.dockerManager.ServiceForceUpdate(ctx, service.ID)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
