package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentAppNotificationSettingsVersion = 1
)

type AppNotificationSettings struct {
	Deployment *AppDeploymentNtfnSettings `json:"deployment,omitempty"`
}

type AppDeploymentNtfnSettings struct {
	Success *AppDeploymentTargetNtfnSettings `json:"success,omitempty"`
	Failure *AppDeploymentTargetNtfnSettings `json:"failure,omitempty"`
}

func (s *AppDeploymentNtfnSettings) HasViaEmailNtfnSettings() bool {
	return (s.Success != nil && s.Success.ViaEmail != nil) || (s.Failure != nil && s.Failure.ViaEmail != nil)
}

func (s *AppDeploymentNtfnSettings) HasViaSlackNtfnSettings() bool {
	return (s.Success != nil && s.Success.ViaSlack != nil) || (s.Failure != nil && s.Failure.ViaSlack != nil)
}

func (s *AppDeploymentNtfnSettings) HasViaDiscordNtfnSettings() bool {
	return (s.Success != nil && s.Success.ViaDiscord != nil) || (s.Failure != nil && s.Failure.ViaDiscord != nil)
}

type AppDeploymentTargetNtfnSettings struct {
	ViaEmail   *EmailNtfnSetting   `json:"viaEmail,omitempty"`
	ViaSlack   *SlackNtfnSetting   `json:"viaSlack,omitempty"`
	ViaDiscord *DiscordNtfnSetting `json:"viaDiscord,omitempty"`
}

func (s *AppNotificationSettings) GetType() base.SettingType {
	return base.SettingTypeAppNotification
}

func (s *AppNotificationSettings) GetRefSettingIDs() []string {
	res := make([]string, 0, 5) //nolint
	res = append(res, s.GetRefEmailIDs()...)
	res = append(res, s.GetRefIMServiceIDs()...)
	return res
}

func (s *AppNotificationSettings) GetRefEmailIDs() (res []string) {
	if s.Deployment != nil {
		if s.Deployment.Success != nil {
			res = append(res, s.Deployment.Success.ViaEmail.GetRefSettingIDs()...)
		}
		if s.Deployment.Failure != nil {
			res = append(res, s.Deployment.Failure.ViaEmail.GetRefSettingIDs()...)
		}
	}
	res = gofn.ToSet(res)
	return
}

func (s *AppNotificationSettings) GetRefIMServiceIDs() (res []string) {
	if s.Deployment != nil {
		if s.Deployment.Success != nil {
			res = append(res, s.Deployment.Success.ViaSlack.GetRefSettingIDs()...)
			res = append(res, s.Deployment.Success.ViaDiscord.GetRefSettingIDs()...)
		}
		if s.Deployment.Failure != nil {
			res = append(res, s.Deployment.Failure.ViaSlack.GetRefSettingIDs()...)
			res = append(res, s.Deployment.Failure.ViaDiscord.GetRefSettingIDs()...)
		}
	}
	res = gofn.ToSet(res)
	return
}

func (s *AppNotificationSettings) HasDeploymentNotificationSettings() bool {
	return s.Deployment != nil && (s.Deployment.Success != nil || s.Deployment.Failure != nil)
}

func (s *Setting) AsAppNotificationSettings() (*AppNotificationSettings, error) {
	return parseSettingAs(s, func() *AppNotificationSettings { return &AppNotificationSettings{} })
}

func (s *Setting) MustAsAppNotificationSettings() *AppNotificationSettings {
	return gofn.Must(s.AsAppNotificationSettings())
}
