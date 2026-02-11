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

	// Cron job notification
	EmailSendCronTaskNotification(ctx context.Context, db database.IDB,
		data *EmailMsgDataCronTaskNotification) error
	SlackSendCronTaskNotification(ctx context.Context, db database.IDB,
		data *SlackMsgDataCronTaskNotification) error
	DiscordSendCronTaskNotification(ctx context.Context, db database.IDB,
		data *DiscordMsgDataCronTaskNotification) error

	// Healthcheck
	EmailSendHealthcheckNotification(ctx context.Context, db database.IDB,
		data *EmailMsgDataHealthcheckNotification) error
	SlackSendHealthcheckNotification(ctx context.Context, db database.IDB,
		data *SlackMsgDataHealthcheckNotification) error
	DiscordSendHealthcheckNotification(ctx context.Context, db database.IDB,
		data *DiscordMsgDataHealthcheckNotification) error
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
