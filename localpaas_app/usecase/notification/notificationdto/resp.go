package notificationdto

import (
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DefaultResultNtfnSettingResp struct {
	Success *DefaultTargetNtfnSettingResp `json:"success"`
	Failure *DefaultTargetNtfnSettingResp `json:"failure"`
}

type DefaultTargetNtfnSettingResp struct {
	ViaEmail   *EmailNtfnSettingResp   `json:"viaEmail"`
	ViaSlack   *SlackNtfnSettingResp   `json:"viaSlack"`
	ViaDiscord *DiscordNtfnSettingResp `json:"viaDiscord"`
}

type EmailNtfnSettingResp struct {
	Sender           *basedto.NamedObjectResp `json:"sender"`
	ToProjectMembers bool                     `json:"toProjectMembers"`
	ToProjectOwners  bool                     `json:"toProjectOwners"`
	ToAllAdmins      bool                     `json:"toAllAdmins"`
	ToAddresses      []string                 `json:"toAddresses"`
}

type SlackNtfnSettingResp struct {
	Webhook *basedto.NamedObjectResp `json:"webhook"`
}

type DiscordNtfnSettingResp struct {
	Webhook *basedto.NamedObjectResp `json:"webhook"`
}
