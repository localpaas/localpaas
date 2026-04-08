package traefikserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (s *service) ReloadTraefikConfig(ctx context.Context, restartServiceOnFailure bool) error {
	// Traefik automatically watches the dynamic configuration directory and reloads changes.

	err := s.reloadTraefikConfig(ctx)
	if err == nil {
		return nil
	}
	if !restartServiceOnFailure {
		return apperrors.Wrap(err)
	}
	err = s.RestartTraefikSwarmService(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *service) reloadTraefikConfig(ctx context.Context) error {
	service, err := s.dockerManager.ServiceGetByName(ctx, traefikServiceName)
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
	if len(containerIDs) == 0 {
		return apperrors.NewNotFound("Traefik service")
	}

	errMap := s.dockerManager.ContainerKillMulti(ctx, containerIDs, "SIGHUP")
	for _, err := range errMap {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *service) ResetTraefikConfig(ctx context.Context) error {
	// Since Traefik dynamic configuration is file-based and managed per-app,
	// there is no master "nginx.conf" to template and overwrite.
	// If the global dynamic_conf.yml needs regeneration, it would happen here.
	// For now, we assume global config is managed by the deployment scripts.
	return nil
}
