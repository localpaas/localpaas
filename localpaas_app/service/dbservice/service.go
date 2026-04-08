package dbservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	MigrateData(ctx context.Context, db database.IDB) error
}
