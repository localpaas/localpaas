package webhookdto

import (
	"net/http"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type HandleWebhookGithubReq struct {
	Request *http.Request `json:"-"`
}

func NewHandleWebhookGithubReq() *HandleWebhookGithubReq {
	return &HandleWebhookGithubReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *HandleWebhookGithubReq) Validate() apperrors.ValidationErrors {
	return nil
}

type HandleWebhookGithubResp struct {
	Meta *basedto.Meta `json:"meta"`
}
