package docker

import (
	"context"

	"github.com/docker/docker/api/types/registry"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (m *Manager) RegistryLogin(
	ctx context.Context,
	auth *registry.AuthConfig,
) (*registry.AuthenticateOKBody, error) {
	resp, err := m.client.RegistryLogin(ctx, *auth)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}
