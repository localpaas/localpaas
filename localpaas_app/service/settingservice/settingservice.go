package settingservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
)

type SettingService interface {
	PersistSettingData(ctx context.Context, db database.IDB, data *PersistingSettingData) error

	LoadReferenceSettings(ctx context.Context, db database.IDB, project *entity.Project, app *entity.App,
		requireActive bool, inSettings ...*entity.Setting) (settingMap map[string]*entity.Setting, err error)

	// Events
	OnCreate(ctx context.Context, db database.IDB, event *CreateEvent) error
	OnUpdate(ctx context.Context, db database.IDB, event *UpdateEvent) error
	OnDelete(ctx context.Context, db database.IDB, event *DeleteEvent) error
}

func NewSettingService(
	settingRepo repository.SettingRepo,
	healthcheckSettingsRepo cacherepository.HealthcheckSettingsRepo,
	permissionManager permission.Manager,
) SettingService {
	return &settingService{
		settingRepo:             settingRepo,
		healthcheckSettingsRepo: healthcheckSettingsRepo,
		permissionManager:       permissionManager,
	}
}

type settingService struct {
	settingRepo             repository.SettingRepo
	healthcheckSettingsRepo cacherepository.HealthcheckSettingsRepo
	permissionManager       permission.Manager
}
