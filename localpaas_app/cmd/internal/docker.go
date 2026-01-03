package internal

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/services/docker"
)

func InitDockerManager(lc fx.Lifecycle, manager *docker.Manager, logger logging.Logger) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("initializing docker manager ...")

			resp, err := manager.ContainerExec(ctx, "e14811a15d8bd9dc420f46414d0857cbd56d5b9ba8f7f40c44afe3bf1333c2f9",
				&container.ExecOptions{
					Cmd:          []string{"rm", "/bin"},
					AttachStdout: true,
					AttachStderr: true,
				})
			if err != nil {
				print(err)
			}

			logChan, _ := docker.StartScanningLog(ctx, io.NopCloser(resp.Reader), batchrecvchan.Options{})
			for msgs := range logChan {
				for _, msg := range msgs {
					print("XXXXXXXXXXXXXXXXXXX ", msg.Type, " ", msg.Data)
				}
			}
			resp.Close()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing docker manager ...")
			return manager.Close()
		},
	})
	return nil
}
