package webhookdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	webhookNameMaxLen   = 100
	webhookSecretMaxLen = 200
)

type CreateWebhookReq struct {
	settings.CreateSettingReq
	*WebhookBaseReq
}

type WebhookBaseReq struct {
	Name   string           `json:"name"`
	Kind   base.WebhookKind `json:"kind"`
	Secret string           `json:"secret"`
}

func (req *WebhookBaseReq) ToEntity() *entity.Webhook {
	return &entity.Webhook{
		Kind:   req.Kind,
		Secret: req.Secret,
	}
}

func (req *WebhookBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Name, true, 1, webhookNameMaxLen, field+"name")...)
	res = append(res, basedto.ValidateStrIn(&req.Kind, true, base.AllWebhookKinds, field+"kind")...)
	res = append(res, basedto.ValidateStr(&req.Secret, true, 1, webhookSecretMaxLen, field+"secret")...)
	return res
}

func NewCreateWebhookReq() *CreateWebhookReq {
	return &CreateWebhookReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateWebhookResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
