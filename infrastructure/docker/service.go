package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/pkg/tracerr"
)

type ServiceListOption func(options *swarm.ServiceListOptions)

func (m *Manager) ServiceList(ctx context.Context, options ...ServiceListOption) ([]swarm.Service, error) {
	opts := swarm.ServiceListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ServiceList(ctx, opts)
	if err != nil {
		return nil, tracerr.Wrap(err, "error listing service")
	}
	return resp, nil
}

type ServiceCreateOption func(options *swarm.ServiceCreateOptions)

func (m *Manager) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options ...ServiceCreateOption) (
	*swarm.ServiceCreateResponse, error) {
	opts := swarm.ServiceCreateOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ServiceCreate(ctx, service, opts)
	if err != nil {
		return nil, tracerr.Wrap(err, "error creating service")
	}
	return &resp, nil
}

func (m *Manager) ServiceRemove(ctx context.Context, serviceID string) error {
	err := m.client.ServiceRemove(ctx, serviceID)
	if err != nil {
		return tracerr.Wrap(err, "error removing service")
	}
	return nil
}

type ServiceInspectOption func(*swarm.ServiceInspectOptions)

func (m *Manager) ServiceInspect(ctx context.Context, serviceID string, options ...ServiceInspectOption) (
	*swarm.Service, []byte, error) {
	opts := swarm.ServiceInspectOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, data, err := m.client.ServiceInspectWithRaw(ctx, serviceID, opts)
	if err != nil {
		return nil, nil, tracerr.Wrap(err, "error inspecting service")
	}
	return &resp, data, nil
}

func (m *Manager) ServiceExists(ctx context.Context, serviceID string) bool {
	resp, _, err := m.ServiceInspect(ctx, serviceID)
	return err == nil && resp != nil
}

type ContainerLogsOption func(*container.LogsOptions)

func (m *Manager) ServiceLogs(ctx context.Context, serviceID string, options ...ContainerLogsOption) (
	io.ReadCloser, error) {
	opts := container.LogsOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ServiceLogs(ctx, serviceID, opts)
	if err != nil {
		return nil, tracerr.Wrap(err, "error getting service logs")
	}
	return resp, nil
}
