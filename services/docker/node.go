package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

type NodeListOption func(*swarm.NodeListOptions)

func (m *Manager) NodeList(ctx context.Context, options ...NodeListOption) ([]swarm.Node, error) {
	opts := swarm.NodeListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NodeList(ctx, opts)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return resp, nil
}

func (m *Manager) NodeInspect(ctx context.Context, nodeID string) (*swarm.Node, []byte, error) {
	resp, data, err := m.client.NodeInspectWithRaw(ctx, nodeID)
	if err != nil {
		return nil, nil, tracerr.Wrap(err)
	}
	return &resp, data, nil
}

func (m *Manager) NodeUpdate(ctx context.Context, nodeID string, version *swarm.Version, spec *swarm.NodeSpec) error {
	if spec == nil {
		return nil
	}
	if version == nil {
		resp, _, err := m.client.NodeInspectWithRaw(ctx, nodeID)
		if err != nil {
			return tracerr.Wrap(err)
		}
		version = &resp.Version
	}
	err := m.client.NodeUpdate(ctx, nodeID, *version, *spec)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

type NodeRemoveOption func(*swarm.NodeRemoveOptions)

func NodeRemoveForce(force bool) NodeRemoveOption {
	return func(opts *swarm.NodeRemoveOptions) {
		opts.Force = force
	}
}

func (m *Manager) NodeRemove(ctx context.Context, nodeID string, options ...NodeRemoveOption) error {
	opts := swarm.NodeRemoveOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	err := m.client.NodeRemove(ctx, nodeID, opts)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
