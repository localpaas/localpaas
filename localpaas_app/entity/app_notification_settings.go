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

func (s *AppDeploymentNtfnSettings) HasViaEmailNotificationSettings() bool {
	return (s.Success != nil && s.Success.ViaEmail != nil) || (s.Failure != nil && s.Failure.ViaEmail != nil)
}

func (s *AppDeploymentNtfnSettings) HasViaSlackNotificationSettings() bool {
	return (s.Success != nil && s.Success.ViaSlack != nil) || (s.Failure != nil && s.Failure.ViaSlack != nil)
}

func (s *AppDeploymentNtfnSettings) HasViaDiscordNotificationSettings() bool {
	return (s.Success != nil && s.Success.ViaDiscord != nil) || (s.Failure != nil && s.Failure.ViaDiscord != nil)
}

type AppDeploymentTargetNtfnSettings struct {
	ViaEmail   *AppEmailNtfnSettings   `json:"viaEmail,omitempty"`
	ViaSlack   *AppSlackNtfnSettings   `json:"viaSlack,omitempty"`
	ViaDiscord *AppDiscordNtfnSettings `json:"viaDiscord,omitempty"`
}

type AppEmailNtfnSettings struct {
	Sender           ObjectID `json:"sender"`
	ToProjectMembers bool     `json:"toProjectMembers,omitempty"`
	ToProjectOwners  bool     `json:"toProjectOwners,omitempty"`
	ToAllAdmins      bool     `json:"toAllAdmins,omitempty"`
}

type AppSlackNtfnSettings struct {
	Webhook ObjectID `json:"webhook"`
}

type AppDiscordNtfnSettings struct {
	Webhook ObjectID `json:"webhook"`
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
	if s.Deployment != nil && s.Deployment.Success != nil {
		if s.Deployment.Success.ViaEmail != nil {
			res = append(res, s.Deployment.Success.ViaEmail.Sender.ID)
		}
	}
	if s.Deployment != nil && s.Deployment.Failure != nil {
		if s.Deployment.Failure.ViaEmail != nil {
			res = append(res, s.Deployment.Failure.ViaEmail.Sender.ID)
		}
	}
	res = gofn.ToSet(res)
	return
}

func (s *AppNotificationSettings) GetRefIMServiceIDs() (res []string) {
	if s.Deployment != nil && s.Deployment.Success != nil {
		if s.Deployment.Success.ViaSlack != nil {
			res = append(res, s.Deployment.Success.ViaSlack.Webhook.ID)
		}
		if s.Deployment.Success.ViaDiscord != nil {
			res = append(res, s.Deployment.Success.ViaDiscord.Webhook.ID)
		}
	}
	if s.Deployment != nil && s.Deployment.Failure != nil {
		if s.Deployment.Failure.ViaSlack != nil {
			res = append(res, s.Deployment.Failure.ViaSlack.Webhook.ID)
		}
		if s.Deployment.Failure.ViaDiscord != nil {
			res = append(res, s.Deployment.Failure.ViaDiscord.Webhook.ID)
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
