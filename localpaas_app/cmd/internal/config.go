package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
)

var (
	ErrInvalidConfig = errors.New("invalid config")
)

func InitConfig(lc fx.Lifecycle, cfg *config.Config, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := validateConfig(cfg, logger); err != nil {
				return apperrors.Wrap(err)
			}

			// Register the channel to receive SIGHUP signal
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGHUP)

			// Start a goroutine to handle incoming signals
			go func() {
				for {
					sig := <-sigs // Block until a signal is received
					switch sig {
					case syscall.SIGHUP:
						logger.Info("SIGHUP received: Reloading configuration...")
						_, err := config.ReloadConfig()
						if err != nil {
							logger.Errorf("Failed to load configuration: %s", err)
						} else {
							logger.Info("SIGHUP handling: Configuration reloaded.")
						}
					default:
						// Do nothing
					}
				}
			}()
			return nil
		},
	})
}

func validateConfig(cfg *config.Config, logger logging.Logger) error {
	logger.Info("validating app config...")
	isProdEnv := cfg.IsProdEnv()

	// JWT secret must not be empty or a trivial value
	if isProdEnv && len(cfg.Session.JWTSecret) < 10 {
		return fmt.Errorf("%w: invalid JWT secret for production", ErrInvalidConfig)
	}

	// Basic auth password must not be empty or a trivial value
	if isProdEnv && len(cfg.Session.BasicAuthPassword) < 10 {
		return fmt.Errorf("%w: basic auth password is invalid for production", ErrInvalidConfig)
	}

	return nil
}
