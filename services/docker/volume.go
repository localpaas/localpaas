package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/volume"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

type VolumeListOption func(*volume.ListOptions)

func (m *Manager) VolumeList(ctx context.Context, options ...VolumeListOption) (*volume.ListResponse, error) {
	opts := volume.ListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.VolumeList(ctx, opts)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &resp, nil
}

func (m *Manager) VolumeCreate(ctx context.Context, options *volume.CreateOptions) (*volume.Volume, error) {
	resp, err := m.client.VolumeCreate(ctx, *options)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &resp, nil
}

func (m *Manager) VolumeUpdate(ctx context.Context, volumeID string, version *swarm.Version,
	options *volume.UpdateOptions) error {
	if options == nil {
		options = &volume.UpdateOptions{}
	}
	if version == nil {
		resp, _, err := m.client.VolumeInspectWithRaw(ctx, volumeID)
		if err != nil {
			return tracerr.Wrap(err)
		}
		version = &resp.ClusterVolume.Version
	}
	err := m.client.VolumeUpdate(ctx, volumeID, *version, *options)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (m *Manager) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	err := m.client.VolumeRemove(ctx, volumeID, force)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (m *Manager) VolumeInspect(ctx context.Context, volumeID string) (*volume.Volume, []byte, error) {
	resp, raw, err := m.client.VolumeInspectWithRaw(ctx, volumeID)
	if err != nil {
		return nil, nil, tracerr.Wrap(err)
	}
	return &resp, raw, nil
}
