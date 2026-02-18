package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentAppNotificationSettingsVersion = 1
)

var _ = registerSettingParser(base.SettingTypeAppNotification, &appNotificationSettingsParser{})

type appNotificationSettingsParser struct {
}

func (s *appNotificationSettingsParser) New() SettingData {
	return &AppNotificationSettings{}
}

type AppNotificationSettings struct {
	Deployment *DefaultResultNotifSetting `json:"deployment,omitempty"`
}

func (s *AppNotificationSettings) GetType() base.SettingType {
	return base.SettingTypeAppNotification
}

func (s *AppNotificationSettings) GetRefObjectIDs() *RefObjectIDs {
	return s.Deployment.GetRefObjectIDs()
}

func (s *AppNotificationSettings) HasDeploymentNotifSetting() bool {
	return s.Deployment != nil && (s.Deployment.Success != nil || s.Deployment.Failure != nil)
}

func (s *Setting) AsAppNotificationSettings() (*AppNotificationSettings, error) {
	return parseSettingAs[*AppNotificationSettings](s)
}

func (s *Setting) MustAsAppNotificationSettings() *AppNotificationSettings {
	return gofn.Must(s.AsAppNotificationSettings())
}
