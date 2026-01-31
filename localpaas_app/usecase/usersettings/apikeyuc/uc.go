package apikeyuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type APIKeyUC struct {
	db             *database.DB
	settingRepo    repository.SettingRepo
	settingService settingservice.SettingService
}

func NewAPIKeyUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	settingService settingservice.SettingService,
) *APIKeyUC {
	return &APIKeyUC{
		db:             db,
		settingRepo:    settingRepo,
		settingService: settingService,
	}
}
