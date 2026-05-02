package networkservice

import (
	"context"

	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/api/types/swarm"
	"github.com/moby/moby/client"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type Service interface {
	GetProjectNetwork(ctx context.Context, project *entity.Project) (*network.Inspect, error)
	CreateProjectNetwork(ctx context.Context, project *entity.Project) (*client.NetworkCreateResult, error)
	RemoveProjectNetwork(ctx context.Context, project *entity.Project) error

	UpdateAppGlobalRoutingNetwork(ctx context.Context, app *entity.App, service *swarm.Service,
		httpSettings *entity.Setting) error
}
