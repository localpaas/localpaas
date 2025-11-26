package docker

import (
	"github.com/docker/docker/client"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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
		return nil, apperrors.NewInfra(err)
	}
	manager.client = c
	return manager, nil
}

func (m *Manager) Close() error {
	return m.client.Close() //nolint:wrapcheck
}
