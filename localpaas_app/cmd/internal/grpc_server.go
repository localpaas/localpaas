package internal

import (
	"context"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
)

// GrpcRegistrar defines a callback type for registering gRPC services.
type GrpcRegistrar func(s *grpc.Server)

func InitGrpcServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	logger logging.Logger,
	registerFn GrpcRegistrar,
) {
	// gRPC agent listens on port 9090
	port := "9090"
	server := grpc.NewServer()

	// Invoke the register function to register services
	registerFn(server)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			lis, err := net.Listen("tcp", ":"+port)
			if err != nil {
				return apperrors.Wrap(err)
			}
			logger.Infof("gRPC Server listening on port %s ...", port)
			go func() {
				if err := server.Serve(lis); err != nil {
					logger.Errorf("gRPC Server stopped: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping gRPC Server ...")
			server.GracefulStop()
			return nil
		},
	})
}
