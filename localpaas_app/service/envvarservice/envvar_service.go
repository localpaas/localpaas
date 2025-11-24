package envvarservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type EnvVarService interface {
	BuildAppEnv(ctx context.Context, db database.IDB, app *entity.App, buildPhase bool) ([]*EnvVar, error)
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
