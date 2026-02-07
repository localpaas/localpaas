package notificationdto

import (
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type DefaultResultNtfnSettingReq struct {
	Success *DefaultTargetNtfnSettingReq `json:"success"`
	Failure *DefaultTargetNtfnSettingReq `json:"failure"`
}

func (req *DefaultResultNtfnSettingReq) ToEntity() *entity.DefaultResultNtfnSetting {
	return &entity.DefaultResultNtfnSetting{
		Success: req.Success.ToEntity(),
		Failure: req.Failure.ToEntity(),
	}
}

type DefaultTargetNtfnSettingReq struct {
	ViaEmail   *EmailNtfnSettingReq   `json:"viaEmail"`
	ViaSlack   *SlackNtfnSettingReq   `json:"viaSlack"`
	ViaDiscord *DiscordNtfnSettingReq `json:"viaDiscord"`
}

func (req *DefaultTargetNtfnSettingReq) ToEntity() *entity.DefaultTargetNtfnSetting {
	return &entity.DefaultTargetNtfnSetting{
		ViaEmail:   req.ViaEmail.ToEntity(),
		ViaSlack:   req.ViaSlack.ToEntity(),
		ViaDiscord: req.ViaDiscord.ToEntity(),
	}
}

type EmailNtfnSettingReq struct {
	Sender           basedto.ObjectIDReq `json:"sender"`
	ToProjectMembers bool                `json:"toProjectMembers"`
	ToProjectOwners  bool                `json:"toProjectOwners"`
	ToAllAdmins      bool                `json:"toAllAdmins"`
	ToAddresses      []string            `json:"toAddresses"`
}

func (req *EmailNtfnSettingReq) ToEntity() *entity.EmailNtfnSetting {
	if req == nil {
		return nil
	}
	return &entity.EmailNtfnSetting{
		Sender:           entity.ObjectID{ID: req.Sender.ID},
		ToProjectMembers: req.ToProjectMembers,
		ToProjectOwners:  req.ToProjectOwners,
		ToAllAdmins:      req.ToAllAdmins,
		ToAddresses:      req.ToAddresses,
	}
}

type SlackNtfnSettingReq struct {
	Webhook basedto.ObjectIDReq `json:"webhook"`
}

func (req *SlackNtfnSettingReq) ToEntity() *entity.SlackNtfnSetting {
	if req == nil {
		return nil
	}
	return &entity.SlackNtfnSetting{
		Webhook: entity.ObjectID{ID: req.Webhook.ID},
	}
}

type DiscordNtfnSettingReq struct {
	Webhook basedto.ObjectIDReq `json:"webhook"`
}

func (req *DiscordNtfnSettingReq) ToEntity() *entity.DiscordNtfnSetting {
	if req == nil {
		return nil
	}
	return &entity.DiscordNtfnSetting{
		Webhook: entity.ObjectID{ID: req.Webhook.ID},
	}
}
