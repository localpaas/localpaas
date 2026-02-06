package entity

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
