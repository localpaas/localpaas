package awss3dto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateAWSS3MetaReq struct {
	settings.UpdateSettingMetaReq
}

func NewUpdateAWSS3MetaReq() *UpdateAWSS3MetaReq {
	return &UpdateAWSS3MetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAWSS3MetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAWSS3MetaResp struct {
	Meta *basedto.Meta `json:"meta"`
}
