package server

import (
	"context"
	"fmt"

	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	agentproto "github.com/localpaas/localpaas/localpaas_app/interface/agent/proto"
)

// AgentServer implements agentproto.AgentServiceServer
type AgentServer struct {
	agentproto.UnimplementedAgentServiceServer
	logger logging.Logger
}

func NewAgentServer(logger logging.Logger) *AgentServer {
	return &AgentServer{logger: logger}
}

func (s *AgentServer) Ping(ctx context.Context, req *agentproto.PingRequest) (*agentproto.PingResponse, error) {
	s.logger.Infof("Received Ping request with message: %s", req.GetMessage())
	return &agentproto.PingResponse{
		Message: fmt.Sprintf("Pong: %s", req.GetMessage()),
	}, nil
}
