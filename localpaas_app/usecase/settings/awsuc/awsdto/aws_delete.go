package awsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteAWSReq struct {
	settings.DeleteSettingReq
}

func NewDeleteAWSReq() *DeleteAWSReq {
	return &DeleteAWSReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteAWSReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteAWSResp struct {
	Meta *basedto.Meta `json:"meta"`
}
