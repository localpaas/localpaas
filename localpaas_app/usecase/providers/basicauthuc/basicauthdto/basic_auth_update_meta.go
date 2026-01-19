package basicauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type UpdateBasicAuthMetaReq struct {
	providers.UpdateSettingMetaReq
}

func NewUpdateBasicAuthMetaReq() *UpdateBasicAuthMetaReq {
	return &UpdateBasicAuthMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateBasicAuthMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateBasicAuthMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
