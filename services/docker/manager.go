package docker

import (
	"github.com/docker/docker/client"

	"github.com/localpaas/localpaas/pkg/tracerr"
)

type Manager struct {
	client *client.Client
}

func New() (*Manager, error) {
	manager := &Manager{}
	c, err := client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	manager.client = c
	return manager, nil
}

func (m *Manager) Close() error {
	return m.client.Close() //nolint:wrapcheck
}
