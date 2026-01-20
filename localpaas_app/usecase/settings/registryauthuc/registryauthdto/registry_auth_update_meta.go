package registryauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateRegistryAuthMetaReq struct {
	settings.UpdateSettingMetaReq
}

func NewUpdateRegistryAuthMetaReq() *UpdateRegistryAuthMetaReq {
	return &UpdateRegistryAuthMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateRegistryAuthMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateRegistryAuthMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
