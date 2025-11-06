package internal

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
)

var (
	ErrInvalidConfig = errors.New("invalid config")
)

func ValidateConfig(lc fx.Lifecycle, cfg *config.Config, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("validating app config...")
			isProdEnv := cfg.IsProdEnv()

			// JWT secret must not be empty or a trivial value
			if isProdEnv && len(cfg.Session.JWT.Secret) < 10 {
				return fmt.Errorf("%w: invalid JWT secret for production", ErrInvalidConfig)
			}

			// Basic auth password must not be empty or a trivial value
			if isProdEnv && len(cfg.Session.BasicAuth.Password) < 10 {
				return fmt.Errorf("%w: basic auth password is invalid for production", ErrInvalidConfig)
			}

			return nil
		},
	})
}
