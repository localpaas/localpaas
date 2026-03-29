package traefikservice

import (
	"context"
)

func (s *traefikService) ReloadTraefikConfig(ctx context.Context, restartServiceOnFailure bool) error {
	// Traefik automatically watches the dynamic configuration directory and reloads changes.
	// We do not need to send a SIGHUP signal to the container.
	return nil
}

func (s *traefikService) ResetTraefikConfig(ctx context.Context) error {
	// Since Traefik dynamic configuration is file-based and managed per-app,
	// there is no master "nginx.conf" to template and overwrite.
	// If the global dynamic_conf.yml needs regeneration, it would happen here.
	// For now, we assume global config is managed by the deployment scripts.
	return nil
}
