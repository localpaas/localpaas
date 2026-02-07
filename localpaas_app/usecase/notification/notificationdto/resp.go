package notificationdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
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
	Sender           *settings.BaseSettingResp `json:"sender"`
	ToProjectMembers bool                      `json:"toProjectMembers"`
	ToProjectOwners  bool                      `json:"toProjectOwners"`
	ToAllAdmins      bool                      `json:"toAllAdmins"`
	ToAddresses      []string                  `json:"toAddresses"`
}

type SlackNtfnSettingResp struct {
	Webhook *settings.BaseSettingResp `json:"webhook"`
}

type DiscordNtfnSettingResp struct {
	Webhook *settings.BaseSettingResp `json:"webhook"`
}

func TransformDefaultTargetNtfnSetting(
	setting *entity.DefaultTargetNtfnSetting,
	refSettingMap map[string]*entity.Setting,
) (resp *DefaultTargetNtfnSettingResp, err error) {
	if setting == nil {
		return nil, nil
	}
	if err = copier.Copy(&resp, setting); err != nil {
		return nil, apperrors.Wrap(err)
	}
	if setting.ViaEmail != nil {
		itemResp, _ := settings.TransformSettingBase(refSettingMap[setting.ViaEmail.Sender.ID])
		resp.ViaEmail.Sender = itemResp
	}
	if setting.ViaSlack != nil {
		itemResp, _ := settings.TransformSettingBase(refSettingMap[setting.ViaSlack.Webhook.ID])
		resp.ViaSlack.Webhook = itemResp
	}
	if setting.ViaDiscord != nil {
		itemResp, _ := settings.TransformSettingBase(refSettingMap[setting.ViaDiscord.Webhook.ID])
		resp.ViaDiscord.Webhook = itemResp
	}
	return resp, nil
}

func TransformDefaultResultNtfnSetting(
	setting *entity.DefaultResultNtfnSetting,
	refSettingMap map[string]*entity.Setting,
) (resp *DefaultResultNtfnSettingResp, err error) {
	if setting == nil {
		return nil, nil
	}
	resp = &DefaultResultNtfnSettingResp{}
	resp.Success, err = TransformDefaultTargetNtfnSetting(setting.Success, refSettingMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Failure, err = TransformDefaultTargetNtfnSetting(setting.Failure, refSettingMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
