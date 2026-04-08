package clusterservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	PersistClusterData(ctx context.Context, db database.IDB, data *PersistingClusterData) error

	IsMultiNode(ctx context.Context) (bool, error)
}
