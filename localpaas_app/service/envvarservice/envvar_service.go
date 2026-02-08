package envvarservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type EnvVarService interface {
	BuildAppEnvVars(ctx context.Context, db database.IDB, app *entity.App, buildPhase bool) ([]*EnvVar, error)
	ProcessEnvRefs(ctx context.Context, db database.IDB, app *entity.App, envVars []*entity.EnvVar,
		loadEnvVars bool, loadSecrets bool, buildPhase bool) (res []*EnvVar, err error)
}

func NewEnvVarService(
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
) EnvVarService {
	return &envVarService{
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
	}
}

type envVarService struct {
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
}
