package envvarservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	BuildAppEnvVars(ctx context.Context, db database.IDB, app *entity.App, buildPhase bool) ([]*EnvVar, error)
	ProcessEnvRefs(ctx context.Context, db database.IDB, app *entity.App, envVars []*entity.EnvVar,
		loadEnvVars bool, loadSecrets bool, buildPhase bool) (res []*EnvVar, err error)
}
