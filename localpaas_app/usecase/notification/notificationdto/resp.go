package notificationdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DefaultResultNotifSettingResp struct {
	Success *DefaultTargetNotifSettingResp `json:"success"`
	Failure *DefaultTargetNotifSettingResp `json:"failure"`
}

type DefaultTargetNotifSettingResp struct {
	ViaEmail        *EmailNotifSettingResp   `json:"viaEmail"`
	ViaSlack        *SlackNotifSettingResp   `json:"viaSlack"`
	ViaDiscord      *DiscordNotifSettingResp `json:"viaDiscord"`
	MinSendInterval timeutil.Duration        `json:"minSendInterval"`
}

type EmailNotifSettingResp struct {
	Sender           *settings.BaseSettingResp `json:"sender"`
	ToProjectMembers bool                      `json:"toProjectMembers"`
	ToProjectOwners  bool                      `json:"toProjectOwners"`
	ToAllAdmins      bool                      `json:"toAllAdmins"`
	ToAddresses      []string                  `json:"toAddresses"`
}

type SlackNotifSettingResp struct {
	Webhook *settings.BaseSettingResp `json:"webhook"`
}

type DiscordNotifSettingResp struct {
	Webhook *settings.BaseSettingResp `json:"webhook"`
}

func TransformDefaultTargetNotifSetting(
	setting *entity.DefaultTargetNotifSetting,
	refSettingMap map[string]*entity.Setting,
) (resp *DefaultTargetNotifSettingResp, err error) {
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

func TransformDefaultResultNotifSetting(
	setting *entity.DefaultResultNotifSetting,
	refSettingMap map[string]*entity.Setting,
) (resp *DefaultResultNotifSettingResp, err error) {
	if setting == nil {
		return nil, nil
	}
	resp = &DefaultResultNotifSettingResp{}
	resp.Success, err = TransformDefaultTargetNotifSetting(setting.Success, refSettingMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Failure, err = TransformDefaultTargetNotifSetting(setting.Failure, refSettingMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
