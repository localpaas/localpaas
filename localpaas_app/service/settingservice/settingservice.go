package settingservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type SettingService interface {
	PersistSettingData(ctx context.Context, db database.IDB, data *PersistingSettingData) error

	LoadReferenceObjects(ctx context.Context, db database.IDB, scope base.SettingScope,
		objectID string, parentObjectID string, requireActive bool, errorIfUnavail bool,
		inSettings ...*entity.Setting) (*entity.RefObjects, error)
	LoadReferenceObjectsByIDs(ctx context.Context, db database.IDB, scope base.SettingScope,
		objectID string, parentObjectID string, requireActive bool, errorIfUnavail bool,
		refIDs *entity.RefObjectIDs) (*entity.RefObjects, error)

	// Events
	OnCreate(ctx context.Context, db database.IDB, event *CreateEvent) error
	OnUpdate(ctx context.Context, db database.IDB, event *UpdateEvent) error
	OnDelete(ctx context.Context, db database.IDB, event *DeleteEvent) error
}

func NewSettingService(
	settingRepo repository.SettingRepo,
	appRepo repository.AppRepo,
	healthcheckSettingsRepo cacherepository.HealthcheckSettingsRepo,
	userService userservice.UserService,
	permissionManager permission.Manager,
) SettingService {
	return &settingService{
		settingRepo:             settingRepo,
		appRepo:                 appRepo,
		healthcheckSettingsRepo: healthcheckSettingsRepo,
		userService:             userService,
		permissionManager:       permissionManager,
	}
}

type settingService struct {
	settingRepo             repository.SettingRepo
	appRepo                 repository.AppRepo
	healthcheckSettingsRepo cacherepository.HealthcheckSettingsRepo
	userService             userservice.UserService
	permissionManager       permission.Manager
}
