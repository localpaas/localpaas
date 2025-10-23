package settingservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type SettingService interface {
	PersistSettingData(ctx context.Context, db database.IDB, data *PersistingSettingData) error
}

func NewSettingService(
	settingRepo repository.SettingRepo,
	permissionManager permission.Manager,
) SettingService {
	return &settingService{
		settingRepo:       settingRepo,
		permissionManager: permissionManager,
	}
}

type settingService struct {
	settingRepo       repository.SettingRepo
	permissionManager permission.Manager
}
