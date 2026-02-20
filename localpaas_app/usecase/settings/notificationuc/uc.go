package notificationuc

import (
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/services/docker"
)

type NotificationUC struct {
	*settings.BaseSettingUC
	dockerManager docker.Manager
}

func NewNotificationUC(
	baseSettingUC *settings.BaseSettingUC,
	dockerManager docker.Manager,
) *NotificationUC {
	return &NotificationUC{
		BaseSettingUC: baseSettingUC,
		dockerManager: dockerManager,
	}
}
