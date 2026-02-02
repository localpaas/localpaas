package webhookdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateWebhookMetaReq struct {
	settings.UpdateSettingMetaReq
}

func NewUpdateWebhookMetaReq() *UpdateWebhookMetaReq {
	return &UpdateWebhookMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateWebhookMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateWebhookMetaResp struct {
	Meta *basedto.Meta `json:"meta"`
}
