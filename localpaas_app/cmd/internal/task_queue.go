package internal

import (
	"context"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
)

func InitTaskQueue(
	lc fx.Lifecycle,
	cfg *config.Config,
	taskQueue taskqueue.TaskQueue,
	logger logging.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Infof("initializing task queue...")
			err := taskQueue.Start(cfg)
			if err != nil {
				logger.Errorf("failed to initialize task queue: %v", err)
				return apperrors.Wrap(err)
			}
			return nil
		},

		OnStop: func(ctx context.Context) error {
			logger.Info("stopping task queue ...")
			err := taskQueue.Shutdown()
			if err != nil {
				logger.Errorf("failed to stop task queue: %v", err)
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
}
