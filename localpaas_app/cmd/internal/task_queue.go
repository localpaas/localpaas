package internal

import (
	"context"
	"errors"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/infra/taskqueue"
	"github.com/localpaas/localpaas/localpaas_app/tasks"
)

func InitTaskQueue(lc fx.Lifecycle, client rediscache.Client, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Infof("starting task queue server...")
			err := taskqueue.StartServer(client, tasks.InitTaskHandlers, logger)
			if err != nil {
				logger.Errorf("failed to start task queue server: %v", err)
			}

			logger.Infof("starting task queue client...")
			err = taskqueue.StartClient(client)
			if err != nil {
				logger.Errorf("failed to start task queue server: %v", err)
				return apperrors.Wrap(err)
			}

			return nil
		},

		OnStop: func(ctx context.Context) error {
			logger.Info("stopping task queue client ...")
			err1 := taskqueue.StopClient()
			logger.Info("stopping task queue server ...")
			err2 := taskqueue.StopServer()
			return errors.Join(err1, err2)
		},
	})
}
