package gitsourceuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type GitSourceUC struct {
	db             *database.DB
	settingRepo    repository.SettingRepo
	settingService settingservice.SettingService
}

func NewGitSourceUC(
	db *database.DB,
	settingRepo repository.SettingRepo,
	settingService settingservice.SettingService,
) *GitSourceUC {
	return &GitSourceUC{
		db:             db,
		settingRepo:    settingRepo,
		settingService: settingService,
	}
}
