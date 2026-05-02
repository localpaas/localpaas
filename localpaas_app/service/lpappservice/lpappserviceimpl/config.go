package lpappserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

func (s *service) ReloadLpAppConfig(ctx context.Context) error {
	service, err := s.dockerManager.ServiceGetByName(ctx, base.LocalpaasAppServiceName, false)
	if err != nil {
		return apperrors.Wrap(err)
	}

	listResp, err := s.dockerManager.ServiceContainerList(ctx, service.ID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	containers := listResp.Items
	containerIDs := make([]string, 0, len(containers))
	for i := range containers {
		containerIDs = append(containerIDs, containers[i].ID)
	}

	errMap := s.dockerManager.ContainerKillMulti(ctx, containerIDs, "SIGHUP")
	for _, err := range errMap {
		return apperrors.Wrap(err)
	}
	return nil
}
