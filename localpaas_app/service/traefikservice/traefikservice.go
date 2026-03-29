package traefikservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/services/docker"
)

type TraefikService interface {
	GetTraefikSwarmService(ctx context.Context) (*swarm.Service, error)
	RestartTraefikSwarmService(ctx context.Context) error

	ReloadTraefikConfig(ctx context.Context, restartServiceOnFailure bool) error
	ResetTraefikConfig(ctx context.Context) error

	ApplyAppConfig(ctx context.Context, app *entity.App, service *swarm.Service, data *AppConfigData) error
	RemoveAppConfig(ctx context.Context, app *entity.App, service *swarm.Service) error
}

func NewTraefikService(
	settingRepo repository.SettingRepo,
	dockerManager docker.Manager,
) TraefikService {
	return &traefikService{
		settingRepo:   settingRepo,
		dockerManager: dockerManager,
	}
}

type traefikService struct {
	settingRepo   repository.SettingRepo
	dockerManager docker.Manager
}
