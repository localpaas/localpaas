package main

import (
	"context"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/localpaas/localpaas/localpaas_app/cmd/internal"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	agentproto "github.com/localpaas/localpaas/localpaas_app/interface/agent/proto"
	agentserver "github.com/localpaas/localpaas/localpaas_app/interface/agent/server"
)

const (
	startTimeoutDefault = 15 * time.Second
)

func main() {
	provides := []any{
		context.Background,
		config.LoadConfig,
		logging.NewZapLogger,
		agentserver.NewAgentServer,
		func(agentSrv *agentserver.AgentServer) internal.GrpcRegistrar {
			return func(s *grpc.Server) {
				agentproto.RegisterAgentServiceServer(s, agentSrv)
			}
		},
	}

	app := fx.New(
		fx.StartTimeout(startTimeoutDefault),
		fx.Provide(provides...),
		fx.Invoke(internal.InitLogger),
		fx.Invoke(internal.InitConfig),
		fx.Invoke(internal.InitGrpcServer),
	)

	app.Run()
}
