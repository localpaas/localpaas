package awsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateAWSReq struct {
	settings.UpdateSettingReq
	*AWSBaseReq
}

func NewUpdateAWSReq() *UpdateAWSReq {
	return &UpdateAWSReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAWSReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAWSResp struct {
	Meta *basedto.Meta `json:"meta"`
}
