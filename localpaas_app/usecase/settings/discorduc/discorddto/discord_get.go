package discorddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedWebhook = "****************"
)

type GetDiscordReq struct {
	settings.GetSettingReq
}

func NewGetDiscordReq() *GetDiscordReq {
	return &GetDiscordReq{}
}

func (req *GetDiscordReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetDiscordResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *DiscordResp      `json:"data"`
}

type DiscordResp struct {
	*settings.BaseSettingResp
	Webhook   string `json:"webhook"`
	Encrypted bool   `json:"encrypted,omitempty"`
}

func (resp *DiscordResp) CopyWebhook(field entity.EncryptedField) error {
	resp.Webhook = field.String()
	return nil
}

func TransformDiscord(setting *entity.Setting, objectID string) (resp *DiscordResp, err error) {
	config := setting.MustAsDiscord()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.Webhook.IsEncrypted()
	if resp.Encrypted {
		resp.Webhook = maskedWebhook
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting, objectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
