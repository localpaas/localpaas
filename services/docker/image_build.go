package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/build"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type ImageBuildOption func(options *build.ImageBuildOptions)

func (m *Manager) ImageBuild(
	ctx context.Context,
	buildContext io.Reader,
	options ...ImageBuildOption,
) (*build.ImageBuildResponse, error) {
	opts := build.ImageBuildOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	resp, err := m.client.ImageBuild(ctx, buildContext, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}

func (m *Manager) ImageBuildCancel(
	ctx context.Context,
	buildID string,
) error {
	err := m.client.BuildCancel(ctx, buildID)
	if err != nil {
		return apperrors.NewInfra(err)
	}
	return nil
}
