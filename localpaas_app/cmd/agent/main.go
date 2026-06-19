package main

import (
	"context"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/localpaas/localpaas/localpaas_app/cmd/internal"
	agentproto "github.com/localpaas/localpaas/localpaas_app/interface/agent/proto"
	agentserver "github.com/localpaas/localpaas/localpaas_app/interface/agent/server"
	"github.com/localpaas/localpaas/localpaas_app/registry"
)

const (
	startTimeoutDefault = 60 * time.Second
)

func main() {
	provides := []any{
		context.Background,
		func(agentSrv *agentserver.AgentServer) internal.GrpcRegistrar {
			return func(s *grpc.Server) {
				agentproto.RegisterAgentServiceServer(s, agentSrv)
				agentproto.RegisterContainerServiceServer(s, agentSrv)
			}
		},
	}
	provides = append(provides, registry.Provides...)

	app := fx.New(
		fx.StartTimeout(startTimeoutDefault),
		fx.Provide(provides...),
		fx.Invoke(internal.InitLogger),
		fx.Invoke(internal.InitConfig),
		fx.Invoke(internal.InitDBConnection),
		fx.Invoke(internal.InitCache),
		fx.Invoke(internal.InitDockerManager),
		fx.Invoke(internal.InitGrpcServer),
	)

	app.Run()
}
