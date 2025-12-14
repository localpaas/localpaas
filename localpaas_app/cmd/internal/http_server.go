package internal

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/server"
)

func InitHTTPServer(lc fx.Lifecycle, cfg *config.Config, srv server.Server, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if isHttpServerEnabled(cfg) {
				go func() {
					logger.Info("starting HTTP server ...")
					if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						logger.Fatalf("start server error: %v", err.Error())
					}
				}()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if isHttpServerEnabled(cfg) {
				logger.Info("stopping HTTP server ...")
				return srv.Stop(ctx)
			}
			return nil
		},
	})
}

func isHttpServerEnabled(cfg *config.Config) bool {
	return cfg.RunMode == config.RunModeApp || cfg.RunMode == config.RunModeEmbeddedWorker
}
