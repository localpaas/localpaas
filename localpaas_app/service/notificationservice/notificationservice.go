package notificationservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type NotificationService interface {
	// App deployment notification
	EmailSendAppDeploymentNotification(ctx context.Context, db database.IDB,
		data *EmailMsgDataAppDeploymentNotification) error
	SlackSendAppDeploymentNotification(ctx context.Context, db database.IDB,
		data *SlackMsgDataAppDeploymentNotification) error
	DiscordSendAppDeploymentNotification(ctx context.Context, db database.IDB,
		data *DiscordMsgDataAppDeploymentNotification) error
}

func NewNotificationService(
	settingRepo repository.SettingRepo,
) NotificationService {
	return &notificationService{
		settingRepo: settingRepo,
	}
}

type notificationService struct {
	settingRepo repository.SettingRepo
}
