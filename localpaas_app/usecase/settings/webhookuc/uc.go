package webhookuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type WebhookUC struct {
	db                       *database.DB
	settingRepo              repository.SettingRepo
	projectSharedSettingRepo repository.ProjectSharedSettingRepo
	settingService           settingservice.SettingService
}

func NewWebhookUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	projectSharedSettingRepo repository.ProjectSharedSettingRepo,
	settingService settingservice.SettingService,
) *WebhookUC {
	return &WebhookUC{
		db:                       db,
		settingRepo:              settingRepo,
		projectSharedSettingRepo: projectSharedSettingRepo,
		settingService:           settingService,
	}
}
