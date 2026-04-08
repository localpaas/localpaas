package lpappservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
)

type Service interface {
	GetLpAppSwarmService(ctx context.Context) (*swarm.Service, error)
	GetLpAppTasks(ctx context.Context) ([]swarm.Task, error)
	RestartLpAppSwarmService(ctx context.Context) error
	ReloadLpAppConfig(ctx context.Context) error

	GetLpDbSwarmService(ctx context.Context) (*swarm.Service, error)
	RestartLpDbSwarmService(ctx context.Context) error

	GetLpCacheSwarmService(ctx context.Context) (*swarm.Service, error)
	RestartLpCacheSwarmService(ctx context.Context) error
}
