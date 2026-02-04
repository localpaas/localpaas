package webhookdto

import (
	"net/http"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	secretMaxLen = 100
)

type HandleRepoWebhookReq struct {
	Request     *http.Request    `json:"-"`
	WebhookKind base.WebhookKind `json:"-"`
	Secret      string           `json:"-"`
}

func NewHandleRepoWebhookReq() *HandleRepoWebhookReq {
	return &HandleRepoWebhookReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *HandleRepoWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(&req.WebhookKind, true,
		base.AllWebhookKinds, "webhookKind")...)
	validators = append(validators, basedto.ValidateStr(&req.Secret, true,
		1, secretMaxLen, "secret")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type HandleRepoWebhookResp struct {
	Meta *basedto.Meta `json:"meta"`
}
