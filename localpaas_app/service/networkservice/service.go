package networkservice

import (
	"context"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type Service interface {
	GetProjectNetwork(ctx context.Context, project *entity.Project) (*network.Inspect, error)
	CreateProjectNetwork(ctx context.Context, project *entity.Project) (*network.CreateResponse, error)
	RemoveProjectNetwork(ctx context.Context, project *entity.Project) error

	UpdateAppGlobalRoutingNetwork(ctx context.Context, app *entity.App, service *swarm.Service,
		httpSettings *entity.Setting) error
}
