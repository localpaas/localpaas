package apikeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateAPIKeyMetaReq struct {
	settings.UpdateSettingMetaReq
}

func NewUpdateAPIKeyMetaReq() *UpdateAPIKeyMetaReq {
	return &UpdateAPIKeyMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAPIKeyMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAPIKeyMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
