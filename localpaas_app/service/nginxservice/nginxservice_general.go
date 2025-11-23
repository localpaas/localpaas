package nginxservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	nginxServiceName = "localpaas_nginx"
)

func (s *nginxService) GetNginxSwarmService(ctx context.Context) (*swarm.Service, error) {
	service, err := s.dockerManager.ServiceGetByName(ctx, nginxServiceName)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return service, nil
}

func (s *nginxService) ReloadNginxConfig(ctx context.Context) error {
	service, err := s.dockerManager.ServiceGetByName(ctx, nginxServiceName)
	if err != nil {
		return apperrors.Wrap(err)
	}

	containers, err := s.dockerManager.ServiceContainerList(ctx, service.ID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	containerIDs := make([]string, 0, len(containers))
	for _, container := range containers {
		containerIDs = append(containerIDs, container.ID)
	}

	errMap := s.dockerManager.ContainerKillMulti(ctx, containerIDs, "SIGHUP")
	for _, err := range errMap {
		return apperrors.Wrap(err)
	}
	return nil
}
