package webhookdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteWebhookReq struct {
	settings.DeleteSettingReq
}

func NewDeleteWebhookReq() *DeleteWebhookReq {
	return &DeleteWebhookReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteWebhookResp struct {
	Meta *basedto.Meta `json:"meta"`
}
