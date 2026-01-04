package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (m *Manager) ContainerExec(
	ctx context.Context,
	containerID string,
	options *container.ExecOptions,
) (string, *types.HijackedResponse, error) {
	_, err := m.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", nil, apperrors.NewInfra(err)
	}

	resp, err := m.client.ContainerExecCreate(ctx, containerID, *options)
	if err != nil {
		return "", nil, apperrors.NewInfra(err)
	}
	execID := resp.ID
	if execID == "" {
		return "", nil, apperrors.New(apperrors.ErrInfraInternal)
	}

	hijackResp, err := m.client.ContainerExecAttach(ctx, execID, container.ExecAttachOptions{
		Detach:      false,
		Tty:         options.Tty,
		ConsoleSize: options.ConsoleSize,
	})
	if err != nil {
		return "", nil, apperrors.NewInfra(err)
	}

	err = m.client.ContainerExecStart(ctx, execID, container.ExecStartOptions{
		Detach:      options.Detach, //nolint
		Tty:         options.Tty,
		ConsoleSize: options.ConsoleSize,
	})
	if err != nil {
		return "", nil, apperrors.NewInfra(err)
	}

	return execID, &hijackResp, nil
}

func (m *Manager) ContainerExecInspect(
	ctx context.Context,
	execID string,
) (*container.ExecInspect, error) {
	resp, err := m.client.ContainerExecInspect(ctx, execID)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return &resp, nil
}
