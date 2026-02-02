package webhookdto

import (
	"net/http"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type HandleGitWebhookReq struct {
	Request   *http.Request  `json:"-"`
	GitSource base.GitSource `json:"-"`
}

func NewHandleGitWebhookReq() *HandleGitWebhookReq {
	return &HandleGitWebhookReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *HandleGitWebhookReq) Validate() apperrors.ValidationErrors {
	return nil
}

type HandleGitWebhookResp struct {
	Meta *basedto.Meta `json:"meta"`
}
