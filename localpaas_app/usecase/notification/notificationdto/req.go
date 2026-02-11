package notificationdto

import (
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

type DefaultResultNotifSettingReq struct {
	Success *DefaultTargetNotifSettingReq `json:"success"`
	Failure *DefaultTargetNotifSettingReq `json:"failure"`
}

func (req *DefaultResultNotifSettingReq) ToEntity() *entity.DefaultResultNotifSetting {
	return &entity.DefaultResultNotifSetting{
		Success: req.Success.ToEntity(),
		Failure: req.Failure.ToEntity(),
	}
}

type DefaultTargetNotifSettingReq struct {
	ViaEmail        *EmailNotifSettingReq   `json:"viaEmail"`
	ViaSlack        *SlackNotifSettingReq   `json:"viaSlack"`
	ViaDiscord      *DiscordNotifSettingReq `json:"viaDiscord"`
	MinSendInterval timeutil.Duration       `json:"minSendInterval"`
}

func (req *DefaultTargetNotifSettingReq) ToEntity() *entity.DefaultTargetNotifSetting {
	return &entity.DefaultTargetNotifSetting{
		ViaEmail:        req.ViaEmail.ToEntity(),
		ViaSlack:        req.ViaSlack.ToEntity(),
		ViaDiscord:      req.ViaDiscord.ToEntity(),
		MinSendInterval: req.MinSendInterval,
	}
}

type EmailNotifSettingReq struct {
	Sender           basedto.ObjectIDReq `json:"sender"`
	ToProjectMembers bool                `json:"toProjectMembers"`
	ToProjectOwners  bool                `json:"toProjectOwners"`
	ToAllAdmins      bool                `json:"toAllAdmins"`
	ToAddresses      []string            `json:"toAddresses"`
}

func (req *EmailNotifSettingReq) ToEntity() *entity.EmailNotifSetting {
	if req == nil {
		return nil
	}
	return &entity.EmailNotifSetting{
		Sender:           entity.ObjectID{ID: req.Sender.ID},
		ToProjectMembers: req.ToProjectMembers,
		ToProjectOwners:  req.ToProjectOwners,
		ToAllAdmins:      req.ToAllAdmins,
		ToAddresses:      req.ToAddresses,
	}
}

type SlackNotifSettingReq struct {
	Webhook basedto.ObjectIDReq `json:"webhook"`
}

func (req *SlackNotifSettingReq) ToEntity() *entity.SlackNotifSetting {
	if req == nil {
		return nil
	}
	return &entity.SlackNotifSetting{
		Webhook: entity.ObjectID{ID: req.Webhook.ID},
	}
}

type DiscordNotifSettingReq struct {
	Webhook basedto.ObjectIDReq `json:"webhook"`
}

func (req *DiscordNotifSettingReq) ToEntity() *entity.DiscordNotifSetting {
	if req == nil {
		return nil
	}
	return &entity.DiscordNotifSetting{
		Webhook: entity.ObjectID{ID: req.Webhook.ID},
	}
}
