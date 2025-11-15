package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

func (m *Manager) SwarmInspect(ctx context.Context) (*swarm.Swarm, error) {
	resp, err := m.client.SwarmInspect(ctx)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &resp, nil
}
