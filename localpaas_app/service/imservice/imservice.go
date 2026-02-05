package imservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type IMService interface {
	// Slack
	SlackSendAppDeploymentNotification(ctx context.Context, db database.IDB,
		data *SlackMsgDataAppDeploymentNotification) error

	// Discord
	DiscordSendAppDeploymentNotification(ctx context.Context, db database.IDB,
		data *DiscordMsgDataAppDeploymentNotification) error
}

func NewIMService(
	settingRepo repository.SettingRepo,
) IMService {
	return &imService{
		settingRepo: settingRepo,
	}
}

type imService struct {
	settingRepo repository.SettingRepo
}
