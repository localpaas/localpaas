package lpappservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (s *lpAppService) ReloadLpAppConfig(ctx context.Context) error {
	service, err := s.dockerManager.ServiceGetByName(ctx, lpAppServiceName)
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
