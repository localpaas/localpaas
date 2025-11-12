package internal

import (
	"context"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/services/docker"
)

func InitDockerManager(lc fx.Lifecycle, manager *docker.Manager, logger logging.Logger) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("initializing docker manager ...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing docker manager ...")
			return manager.Close()
		},
	})
	return nil
}
