package docker

import (
	"context"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type ContainerListOption func(*container.ListOptions)

func (m *Manager) ContainerList(
	ctx context.Context,
	options ...ContainerListOption,
) ([]container.Summary, error) {
	opts := container.ListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	containers, err := m.client.ContainerList(ctx, opts)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}
	return containers, nil
}

func (m *Manager) ServiceContainerList(
	ctx context.Context,
	serviceID string,
	options ...ContainerListOption,
) ([]container.Summary, error) {
	options = append(options, func(opts *container.ListOptions) {
		opts.All = false
		FilterAdd(&opts.Filters, "label", "com.docker.swarm.service.id="+serviceID)
	})
	return m.ContainerList(ctx, options...)
}

func (m *Manager) ServiceContainerGetActive(
	ctx context.Context,
	serviceID string,
	maxRetry int,
	retryDelay time.Duration,
) (active *container.Summary, all []container.Summary, err error) {
	return m.serviceContainerGetActive(ctx, serviceID, -1, maxRetry, retryDelay)
}

func (m *Manager) serviceContainerGetActive(
	ctx context.Context,
	serviceID string,
	retry int,
	maxRetry int,
	retryDelay time.Duration,
) (active *container.Summary, all []container.Summary, err error) {
	if retry >= maxRetry {
		return nil, nil, nil
	}
	summaries, err := m.ServiceContainerList(ctx, serviceID)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}

	for i := range summaries {
		c := &summaries[i]
		if c.State == container.StateRunning {
			return c, summaries, nil
		}
	}

	time.Sleep(retryDelay)
	return m.serviceContainerGetActive(ctx, serviceID, retry+1, maxRetry, retryDelay)
}

func (m *Manager) ContainerInspect(
	ctx context.Context,
	containerID string,
) (*container.InspectResponse, error) {
	respMap, errMap := m.ContainerInspectMulti(ctx, []string{containerID})
	return respMap[containerID], errMap[containerID]
}

func (m *Manager) ContainerInspectMulti(
	ctx context.Context,
	containerIDs []string,
) (map[string]*container.InspectResponse, map[string]error) {
	if len(containerIDs) == 1 {
		resp, err := m.client.ContainerInspect(ctx, containerIDs[0])
		if err != nil {
			return nil, map[string]error{containerIDs[0]: apperrors.NewInfra(err)}
		}
		return map[string]*container.InspectResponse{containerIDs[0]: &resp}, nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	allResults := make(map[string]*container.InspectResponse, len(containerIDs))
	allErrors := map[string]error{}
	for _, containerID := range containerIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := m.client.ContainerInspect(ctx, containerID)
			mu.Lock()
			if err != nil {
				allErrors[containerID] = apperrors.NewInfra(err)
			} else {
				allResults[containerID] = &resp
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	return allResults, allErrors
}

type ContainerStopOption func(options *container.StopOptions)

func (m *Manager) ContainerRestart(
	ctx context.Context,
	containerID string,
	options ...ContainerStopOption,
) error {
	errMap := m.ContainerRestartMulti(ctx, []string{containerID}, options...)
	for _, err := range errMap {
		return err
	}
	return nil
}

func (m *Manager) ContainerRestartMulti(
	ctx context.Context,
	containerIDs []string,
	options ...ContainerStopOption,
) map[string]error {
	opts := &container.StopOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if len(containerIDs) == 1 {
		err := m.client.ContainerRestart(ctx, containerIDs[0], *opts)
		if err != nil {
			return map[string]error{containerIDs[0]: apperrors.NewInfra(err)}
		}
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	allErrors := map[string]error{}
	for _, containerID := range containerIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := m.client.ContainerRestart(ctx, containerID, *opts)
			if err != nil {
				mu.Lock()
				allErrors[containerID] = apperrors.NewInfra(err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return allErrors
}

func (m *Manager) ContainerKill(
	ctx context.Context,
	containerID string,
	signal string,
) error {
	errMap := m.ContainerKillMulti(ctx, []string{containerID}, signal)
	for _, err := range errMap {
		return err
	}
	return nil
}

func (m *Manager) ContainerKillMulti(
	ctx context.Context,
	containerIDs []string,
	signal string,
) map[string]error {
	if len(containerIDs) == 1 {
		err := m.client.ContainerKill(ctx, containerIDs[0], signal)
		if err != nil {
			return map[string]error{containerIDs[0]: apperrors.NewInfra(err)}
		}
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	allErrors := map[string]error{}
	for _, containerID := range containerIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := m.client.ContainerKill(ctx, containerID, signal)
			if err != nil {
				mu.Lock()
				allErrors[containerID] = apperrors.NewInfra(err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return allErrors
}
