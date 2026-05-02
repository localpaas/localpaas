package docker

import (
	"context"
	"time"

	"github.com/moby/moby/client"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	NetworkDriverOverlay = "overlay"
	NetworkDriverBridge  = "bridge"
)

const (
	NetworkScopeSwarm = "swarm"
	NetworkScopeLocal = "local"
)

type NetworkListOption func(*client.NetworkListOptions)

func (m *manager) NetworkList(
	ctx context.Context,
	options ...NetworkListOption,
) (*client.NetworkListResult, error) {
	opts := client.NetworkListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NetworkList(ctx, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}

type NetworkCreateOption func(*client.NetworkCreateOptions)

func (m *manager) NetworkCreate(
	ctx context.Context,
	name string,
	options ...NetworkCreateOption,
) (*client.NetworkCreateResult, error) {
	opts := client.NetworkCreateOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NetworkCreate(ctx, name, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}

type NetworkRemoveOption func(*client.NetworkRemoveOptions)

func (m *manager) NetworkRemove(
	ctx context.Context,
	idOrName string,
	options ...NetworkRemoveOption,
) (*client.NetworkRemoveResult, error) {
	opts := client.NetworkRemoveOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NetworkRemove(ctx, idOrName, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}

type NetworkInspectOption func(*client.NetworkInspectOptions)

func (m *manager) NetworkInspect(
	ctx context.Context,
	name string,
	options ...NetworkInspectOption,
) (*client.NetworkInspectResult, error) {
	opts := client.NetworkInspectOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NetworkInspect(ctx, name, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}

func (m *manager) NetworkExists(ctx context.Context, name string) bool {
	_, err := m.NetworkInspect(ctx, name)
	return err == nil
}

type NetworkPruneOption func(*client.NetworkPruneOptions)

func (m *manager) NetworkPrune(
	ctx context.Context,
	onlyObjectsOlderThan time.Duration,
	options ...NetworkPruneOption,
) (*client.NetworkPruneResult, error) {
	opts := client.NetworkPruneOptions{}
	if onlyObjectsOlderThan > 0 {
		FilterAdd(&opts.Filters, "until", onlyObjectsOlderThan.String())
	}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.NetworkPrune(ctx, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}
