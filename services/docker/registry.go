package docker

import (
	"context"

	"github.com/docker/docker/api/types/registry"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

func (m *Manager) RegistryLogin(ctx context.Context, auth *registry.AuthConfig) (*registry.AuthenticateOKBody, error) {
	resp, err := m.client.RegistryLogin(ctx, *auth)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &resp, nil
}
