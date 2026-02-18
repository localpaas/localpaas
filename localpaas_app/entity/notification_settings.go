package entity

import "github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"

type DefaultResultNotifSetting struct {
	Success *DefaultTargetNotifSetting `json:"success,omitempty"`
	Failure *DefaultTargetNotifSetting `json:"failure,omitempty"`
}

func (s *DefaultResultNotifSetting) GetRefObjectIDs() *RefObjectIDs {
	res := &RefObjectIDs{}
	if s != nil {
		res.AddRefIDs(s.Success.GetRefObjectIDs())
		res.AddRefIDs(s.Failure.GetRefObjectIDs())
	}
	return res
}

func (s *DefaultResultNotifSetting) HasViaEmailNotifSetting() bool {
	if s == nil {
		return false
	}
	return (s.Success != nil && s.Success.ViaEmail != nil) || (s.Failure != nil && s.Failure.ViaEmail != nil)
}

func (s *DefaultResultNotifSetting) HasViaSlackNotifSetting() bool {
	if s == nil {
		return false
	}
	return (s.Success != nil && s.Success.ViaSlack != nil) || (s.Failure != nil && s.Failure.ViaSlack != nil)
}

func (s *DefaultResultNotifSetting) HasViaDiscordNotifSetting() bool {
	if s == nil {
		return false
	}
	return (s.Success != nil && s.Success.ViaDiscord != nil) || (s.Failure != nil && s.Failure.ViaDiscord != nil)
}

type DefaultTargetNotifSetting struct {
	ViaEmail   *EmailNotifSetting   `json:"viaEmail,omitempty"`
	ViaSlack   *SlackNotifSetting   `json:"viaSlack,omitempty"`
	ViaDiscord *DiscordNotifSetting `json:"viaDiscord,omitempty"`

	MinSendInterval timeutil.Duration `json:"minSendInterval,omitempty"`
}

func (s *DefaultTargetNotifSetting) GetRefObjectIDs() *RefObjectIDs {
	res := &RefObjectIDs{}
	if s != nil {
		res.RefSettingIDs = append(res.RefSettingIDs, s.ViaEmail.GetRefSettingIDs()...)
		res.RefSettingIDs = append(res.RefSettingIDs, s.ViaSlack.GetRefSettingIDs()...)
		res.RefSettingIDs = append(res.RefSettingIDs, s.ViaDiscord.GetRefSettingIDs()...)
	}
	return res
}

func (s *DefaultTargetNotifSetting) HasViaEmailNotifSettings() bool {
	return s.ViaEmail != nil
}

func (s *DefaultTargetNotifSetting) HasViaSlackNotifSettings() bool {
	return s.ViaSlack != nil
}

func (s *DefaultTargetNotifSetting) HasViaDiscordNotifSettings() bool {
	return s.ViaDiscord != nil
}

type EmailNotifSetting struct {
	Sender           ObjectID `json:"sender"`
	ToProjectMembers bool     `json:"toProjectMembers,omitempty"`
	ToProjectOwners  bool     `json:"toProjectOwners,omitempty"`
	ToAllAdmins      bool     `json:"toAllAdmins,omitempty"`
	ToAddresses      []string `json:"toAddresses,omitempty"`
}

func (s *EmailNotifSetting) GetRefSettingIDs() (res []string) {
	if s == nil {
		return res
	}
	res = append(res, s.Sender.ID)
	return res
}

type SlackNotifSetting struct {
	Webhook ObjectID `json:"webhook"`
}

func (s *SlackNotifSetting) GetRefSettingIDs() (res []string) {
	if s == nil {
		return res
	}
	res = append(res, s.Webhook.ID)
	return res
}

type DiscordNotifSetting struct {
	Webhook ObjectID `json:"webhook"`
}

func (s *DiscordNotifSetting) GetRefSettingIDs() (res []string) {
	if s == nil {
		return res
	}
	res = append(res, s.Webhook.ID)
	return res
}
