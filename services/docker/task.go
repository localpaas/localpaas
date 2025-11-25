package docker

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type TaskListOption func(*swarm.TaskListOptions)

func (m *Manager) TaskList(ctx context.Context, options ...TaskListOption) ([]swarm.Task, error) {
	opts := swarm.TaskListOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	tasks, err := m.client.TaskList(ctx, opts)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return tasks, nil
}

func (m *Manager) ServiceTaskList(ctx context.Context, serviceID string, options ...TaskListOption) (
	[]swarm.Task, error) {
	options = append(options, func(opts *swarm.TaskListOptions) {
		if opts.Filters.Len() == 0 {
			opts.Filters = filters.NewArgs()
		}
		opts.Filters.Add("service", serviceID)
	})
	return m.TaskList(ctx, options...)
}
