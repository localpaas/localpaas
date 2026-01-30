package awss3dto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteAWSS3Req struct {
	settings.DeleteSettingReq
}

func NewDeleteAWSS3Req() *DeleteAWSS3Req {
	return &DeleteAWSS3Req{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteAWSS3Req) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteAWSS3Resp struct {
	Meta *basedto.Meta `json:"meta"`
}
