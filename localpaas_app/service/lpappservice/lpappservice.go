package lpappservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/services/docker"
)

type LpAppService interface {
	GetLpAppSwarmService(ctx context.Context) (*swarm.Service, error)
	RestartLpAppSwarmService(ctx context.Context) error
	ReloadLpAppConfig(ctx context.Context) error

	GetLpDbSwarmService(ctx context.Context) (*swarm.Service, error)
	RestartLpDbSwarmService(ctx context.Context) error

	GetLpCacheSwarmService(ctx context.Context) (*swarm.Service, error)
	RestartLpCacheSwarmService(ctx context.Context) error
}

func NewLpAppService(
	settingRepo repository.SettingRepo,
	dockerManager *docker.Manager,
) LpAppService {
	return &lpAppService{
		settingRepo:   settingRepo,
		dockerManager: dockerManager,
	}
}

type lpAppService struct {
	settingRepo   repository.SettingRepo
	dockerManager *docker.Manager
}
