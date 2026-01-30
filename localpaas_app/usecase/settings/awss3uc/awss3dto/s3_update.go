package awss3dto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateAWSS3Req struct {
	settings.UpdateSettingReq
	*AWSS3BaseReq
}

func NewUpdateAWSS3Req() *UpdateAWSS3Req {
	return &UpdateAWSS3Req{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAWSS3Req) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAWSS3Resp struct {
	Meta *basedto.Meta `json:"meta"`
}
