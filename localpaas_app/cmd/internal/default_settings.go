package internal

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

func InitDefaultSettings(
	lc fx.Lifecycle,
	db *database.DB,
	settingService settingservice.SettingService,
	logger logging.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("initializing default settings...")
			if err := settingService.InitDefaults(ctx, db); err != nil {
				return fmt.Errorf("failed to initialize default settings: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
