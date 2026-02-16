package repowebhookdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	webhookSecretMaxLen = 100
)

type CreateRepoWebhookReq struct {
	settings.CreateSettingReq
	*RepoWebhookBaseReq
}

type RepoWebhookBaseReq struct {
	Name   string           `json:"name"`
	Kind   base.WebhookKind `json:"kind"`
	Secret string           `json:"secret"`
}

func (req *RepoWebhookBaseReq) ToEntity() *entity.RepoWebhook {
	return &entity.RepoWebhook{
		Kind:   req.Kind,
		Secret: req.Secret,
	}
}

func (req *RepoWebhookBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Name, true, 1, base.SettingNameMaxLen, field+"name")...)
	res = append(res, basedto.ValidateStrIn(&req.Kind, true, base.AllWebhookKinds, field+"kind")...)
	res = append(res, basedto.ValidateStr(&req.Secret, true, 1, webhookSecretMaxLen, field+"secret")...)
	return res
}

func NewCreateRepoWebhookReq() *CreateRepoWebhookReq {
	return &CreateRepoWebhookReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateRepoWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateRepoWebhookResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
