package imservicedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedWebhook = "****************"
)

type GetIMServiceReq struct {
	settings.GetSettingReq
}

func NewGetIMServiceReq() *GetIMServiceReq {
	return &GetIMServiceReq{}
}

func (req *GetIMServiceReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetIMServiceResp struct {
	Meta *basedto.Meta  `json:"meta"`
	Data *IMServiceResp `json:"data"`
}

type IMServiceResp struct {
	*settings.BaseSettingResp
	Kind      base.IMServiceKind `json:"kind"`
	Slack     *SlackResp         `json:"slack,omitempty"`
	Discord   *DiscordResp       `json:"discord,omitempty"`
	Encrypted bool               `json:"encrypted,omitempty"`
}

type SlackResp struct {
	Webhook string `json:"webhook"`
}

func (resp *SlackResp) CopyWebhook(field entity.EncryptedField) error {
	resp.Webhook = field.String()
	return nil
}

type DiscordResp struct {
	Webhook string `json:"webhook"`
}

func (resp *DiscordResp) CopyWebhook(field entity.EncryptedField) error {
	resp.Webhook = field.String()
	return nil
}

func TransformIMService(setting *entity.Setting) (resp *IMServiceResp, err error) {
	config := setting.MustAsIMService()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Kind = base.IMServiceKind(setting.Kind)

	switch {
	case config.Slack != nil:
		resp.Encrypted = config.Slack.Webhook.IsEncrypted()
		if resp.Encrypted {
			resp.Slack.Webhook = maskedWebhook
		}
	case config.Discord != nil:
		resp.Encrypted = config.Discord.Webhook.IsEncrypted()
		if resp.Encrypted {
			resp.Discord.Webhook = maskedWebhook
		}
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
