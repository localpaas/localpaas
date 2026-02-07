package entity

type DefaultResultNtfnSetting struct {
	Success *DefaultTargetNtfnSetting `json:"success,omitempty"`
	Failure *DefaultTargetNtfnSetting `json:"failure,omitempty"`
}

func (s *DefaultResultNtfnSetting) GetRefSettingIDs() (res []string) {
	if s == nil {
		return res
	}
	res = append(res, s.Success.GetRefSettingIDs()...)
	res = append(res, s.Failure.GetRefSettingIDs()...)
	return res
}

func (s *DefaultResultNtfnSetting) HasViaEmailNtfnSetting() bool {
	return (s.Success != nil && s.Success.ViaEmail != nil) || (s.Failure != nil && s.Failure.ViaEmail != nil)
}

func (s *DefaultResultNtfnSetting) HasViaSlackNtfnSetting() bool {
	return (s.Success != nil && s.Success.ViaSlack != nil) || (s.Failure != nil && s.Failure.ViaSlack != nil)
}

func (s *DefaultResultNtfnSetting) HasViaDiscordNtfnSetting() bool {
	return (s.Success != nil && s.Success.ViaDiscord != nil) || (s.Failure != nil && s.Failure.ViaDiscord != nil)
}

type DefaultTargetNtfnSetting struct {
	ViaEmail   *EmailNtfnSetting   `json:"viaEmail,omitempty"`
	ViaSlack   *SlackNtfnSetting   `json:"viaSlack,omitempty"`
	ViaDiscord *DiscordNtfnSetting `json:"viaDiscord,omitempty"`
}

func (s *DefaultTargetNtfnSetting) GetRefSettingIDs() (res []string) {
	if s == nil {
		return res
	}
	res = append(res, s.ViaEmail.GetRefSettingIDs()...)
	res = append(res, s.ViaSlack.GetRefSettingIDs()...)
	res = append(res, s.ViaDiscord.GetRefSettingIDs()...)
	return res
}

func (s *DefaultTargetNtfnSetting) HasViaEmailNtfnSettings() bool {
	return s.ViaEmail != nil
}

func (s *DefaultTargetNtfnSetting) HasViaSlackNtfnSettings() bool {
	return s.ViaSlack != nil
}

func (s *DefaultTargetNtfnSetting) HasViaDiscordNtfnSettings() bool {
	return s.ViaDiscord != nil
}

type EmailNtfnSetting struct {
	Sender           ObjectID `json:"sender"`
	ToProjectMembers bool     `json:"toProjectMembers,omitempty"`
	ToProjectOwners  bool     `json:"toProjectOwners,omitempty"`
	ToAllAdmins      bool     `json:"toAllAdmins,omitempty"`
	ToAddresses      []string `json:"toAddresses,omitempty"`
}

func (s *EmailNtfnSetting) GetRefSettingIDs() (res []string) {
	if s == nil {
		return res
	}
	res = append(res, s.Sender.ID)
	return res
}

type SlackNtfnSetting struct {
	Webhook ObjectID `json:"webhook"`
}

func (s *SlackNtfnSetting) GetRefSettingIDs() (res []string) {
	if s == nil {
		return res
	}
	res = append(res, s.Webhook.ID)
	return res
}

type DiscordNtfnSetting struct {
	Webhook ObjectID `json:"webhook"`
}

func (s *DiscordNtfnSetting) GetRefSettingIDs() (res []string) {
	if s == nil {
		return res
	}
	res = append(res, s.Webhook.ID)
	return res
}
