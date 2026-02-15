package docker

import (
	"context"

	"github.com/docker/docker/api/types/system"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (m *Manager) SystemInfo(ctx context.Context) (*system.Info, error) {
	resp, err := m.client.Info(ctx)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}
