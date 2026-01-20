package secretdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateSecretMetaReq struct {
	settings.UpdateSettingMetaReq
}

func NewUpdateSecretMetaReq() *UpdateSecretMetaReq {
	return &UpdateSecretMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateSecretMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateSecretMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
