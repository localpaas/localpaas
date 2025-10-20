package internal

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/server"
)

func InitHTTPServer(lc fx.Lifecycle, srv server.Server, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Infof("starting server %v", srv.GetAddress())
			go func() {
				if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatalf("start server error: %v", err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping server ...")
			return srv.Stop(ctx)
		},
	})
}
