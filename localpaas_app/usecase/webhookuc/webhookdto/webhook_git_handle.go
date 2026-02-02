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

type HandleGitWebhookReq struct {
	Request   *http.Request  `json:"-"`
	GitSource base.GitSource `json:"-"`
	Secret    string         `json:"-"`
}

func NewHandleGitWebhookReq() *HandleGitWebhookReq {
	return &HandleGitWebhookReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *HandleGitWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(&req.GitSource, true,
		base.AllGitSources, "gitSource")...)
	validators = append(validators, basedto.ValidateStr(&req.Secret, true,
		1, secretMaxLen, "secret")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type HandleGitWebhookResp struct {
	Meta *basedto.Meta `json:"meta"`
}
