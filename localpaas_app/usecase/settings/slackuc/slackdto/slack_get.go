package slackdto

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

type GetSlackReq struct {
	settings.GetSettingReq
}

func NewGetSlackReq() *GetSlackReq {
	return &GetSlackReq{}
}

func (req *GetSlackReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSlackResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *SlackResp        `json:"data"`
}

type SlackResp struct {
	*settings.BaseSettingResp
	Webhook   string `json:"webhook"`
	Encrypted bool   `json:"encrypted,omitempty"`
}

func (resp *SlackResp) CopyWebhook(field entity.EncryptedField) error {
	resp.Webhook = field.String()
	return nil
}

func TransformSlack(setting *entity.Setting) (resp *SlackResp, err error) {
	config := setting.MustAsSlack()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.Webhook.IsEncrypted()
	if resp.Encrypted {
		resp.Webhook = maskedWebhook
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
