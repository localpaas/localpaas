package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"go.uber.org/fx"

	"github.com/localpaas/localpaas/localpaas_app/cmd/internal"
	"github.com/localpaas/localpaas/localpaas_app/registry"
)

const (
	startTimeoutDefault = 5 * time.Minute
)

func main() {
	provides := []any{
		context.Background,
	}
	provides = append(provides, registry.Provides...)

	app := fx.New(
		fx.StartTimeout(startTimeoutDefault),
		fx.Provide(provides...),
		fx.Invoke(internal.InitLogger),
		fx.Invoke(internal.InitConfig),
		fx.Invoke(internal.InitDBConnection),
		fx.Invoke(internal.MigrateData), // Migrate data structure of JSON columns
		fx.Invoke(internal.InitCache),
		fx.Invoke(internal.InitTaskQueue),
		fx.Invoke(internal.InitDockerManager),
		fx.Invoke(internal.CompleteInstallation),
		fx.Invoke(func() {
			p, _ := os.FindProcess(os.Getpid())
			_ = p.Signal(syscall.SIGTERM)
		}),
	)

	app.Run()
}
