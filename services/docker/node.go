package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/pkg/tracerr"
)

type NodeListOption func(*swarm.NodeListOptions)

func (m *Manager) NodeList(ctx context.Context, options ...NodeListOption) ([]swarm.Node, error) {
	opts := swarm.NodeListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NodeList(ctx, opts)
	if err != nil {
		return nil, tracerr.Wrap(err, "error listing node")
	}
	return resp, nil
}

func (m *Manager) NodeInspect(ctx context.Context, nodeID string) (*swarm.Node, []byte, error) {
	resp, data, err := m.client.NodeInspectWithRaw(ctx, nodeID)
	if err != nil {
		return nil, nil, tracerr.Wrap(err, "error inspecting node")
	}
	return &resp, data, nil
}

type NodeRemoveOption func(*swarm.NodeRemoveOptions)

func (m *Manager) NodeRemove(ctx context.Context, nodeID string, options ...NodeRemoveOption) error {
	opts := swarm.NodeRemoveOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	err := m.client.NodeRemove(ctx, nodeID, opts)
	if err != nil {
		return tracerr.Wrap(err, "error removing node")
	}
	return nil
}
