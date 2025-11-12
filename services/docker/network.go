package docker

import (
	"context"

	"github.com/docker/docker/api/types/network"

	"github.com/localpaas/localpaas/pkg/tracerr"
)

type NetworkListOption func(*network.ListOptions)

func (m *Manager) NetworkList(ctx context.Context, options ...NetworkListOption) ([]network.Summary, error) {
	opts := network.ListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NetworkList(ctx, opts)
	if err != nil {
		return nil, tracerr.Wrap(err, "error listing network")
	}
	return resp, nil
}

type NetworkCreateOption func(*network.CreateOptions)

func (m *Manager) NetworkCreate(ctx context.Context, name string, options ...NetworkCreateOption) (
	*network.CreateResponse, error) {
	opts := network.CreateOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NetworkCreate(ctx, name, opts)
	if err != nil {
		return nil, tracerr.Wrap(err, "error creating network")
	}
	return &resp, nil
}

func (m *Manager) NetworkRemove(ctx context.Context, name string) error {
	err := m.client.NetworkRemove(ctx, name)
	if err != nil {
		return tracerr.Wrap(err, "error removing network")
	}
	return nil
}

type NetworkInspectOption func(*network.InspectOptions)

func (m *Manager) NetworkInspect(ctx context.Context, name string, options ...NetworkInspectOption) (
	*network.Inspect, error) {
	opts := network.InspectOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NetworkInspect(ctx, name, opts)
	if err != nil {
		return nil, tracerr.Wrap(err, "error inspecting network")
	}
	return &resp, nil
}

func (m *Manager) NetworkExists(ctx context.Context, name string) bool {
	_, err := m.NetworkInspect(ctx, name)
	return err == nil
}
