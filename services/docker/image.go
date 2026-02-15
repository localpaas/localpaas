package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/image"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type ImageListOption func(*image.ListOptions)

func (m *Manager) ImageList(
	ctx context.Context,
	options ...ImageListOption,
) ([]image.Summary, error) {
	opts := image.ListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ImageList(ctx, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return resp, nil
}

type ImageCreateOption func(*image.CreateOptions)

func (m *Manager) ImageCreate(
	ctx context.Context,
	name string,
	options ...ImageCreateOption,
) (io.ReadCloser, error) {
	opts := image.CreateOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ImageCreate(ctx, name, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return resp, nil
}

type ImageRemoveOption func(options *image.RemoveOptions)

func (m *Manager) ImageRemove(
	ctx context.Context,
	imageID string,
	options ...ImageRemoveOption,
) ([]image.DeleteResponse, error) {
	opts := image.RemoveOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ImageRemove(ctx, imageID, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return resp, nil
}

func (m *Manager) ImageInspect(
	ctx context.Context,
	imageID string,
) (*image.InspectResponse, error) {
	resp, err := m.client.ImageInspect(ctx, imageID)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}

type ImagePullOption func(options *image.PullOptions)

func (m *Manager) ImagePull(
	ctx context.Context,
	refStr string,
	options ...ImagePullOption,
) (io.ReadCloser, error) {
	opts := image.PullOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ImagePull(ctx, refStr, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return resp, nil
}

type ImagePushOption func(options *image.PushOptions)

func (m *Manager) ImagePush(
	ctx context.Context,
	imageTag string,
	options ...ImagePushOption,
) (io.ReadCloser, error) {
	opts := image.PushOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ImagePush(ctx, imageTag, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return resp, nil
}
