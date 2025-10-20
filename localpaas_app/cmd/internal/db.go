package internal

import (
	"context"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
)

func InitDBConnection(lc fx.Lifecycle, db *database.DB, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("pinging db connection...")
			if err := db.Ping(); err != nil {
				logger.Errorf("failed to use connection %v", err.Error())

				return apperrors.Wrap(err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing db connection...")
			return db.Close()
		},
	})
}
