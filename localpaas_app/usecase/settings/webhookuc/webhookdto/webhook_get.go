package webhookdto

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
	maskedSecret = "****************"
)

type GetWebhookReq struct {
	settings.GetSettingReq
}

func NewGetWebhookReq() *GetWebhookReq {
	return &GetWebhookReq{}
}

func (req *GetWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetWebhookResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data *WebhookResp  `json:"data"`
}

type WebhookResp struct {
	*settings.BaseSettingResp
	Kind      base.WebhookKind `json:"kind"`
	Secret    string           `json:"secret"`
	Encrypted bool             `json:"encrypted,omitempty"`
}

func TransformWebhook(setting *entity.Setting) (resp *WebhookResp, err error) {
	config := setting.MustAsWebhook()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Kind = base.WebhookKind(setting.Kind)

	// resp.Encrypted = config.Secret.IsEncrypted()
	if resp.Encrypted {
		resp.Secret = maskedSecret
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
