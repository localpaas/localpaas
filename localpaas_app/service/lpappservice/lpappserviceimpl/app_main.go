package lpappserviceimpl

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

func (s *service) GetLpAppSwarmService(ctx context.Context) (*swarm.Service, error) {
	service, err := s.dockerManager.ServiceGetByName(ctx, base.LocalpaasAppServiceName)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return service, nil
}

func (s *service) RestartLpAppSwarmService(ctx context.Context) error {
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

func (s *service) GetLpAppTasks(ctx context.Context) ([]swarm.Task, error) {
	service, err := s.GetLpDbSwarmService(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	tasks, err := s.dockerManager.ServiceTaskList(ctx, service.ID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return tasks, nil
}
