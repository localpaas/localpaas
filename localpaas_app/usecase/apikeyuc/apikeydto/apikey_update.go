package apikeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateAPIKeyReq struct {
	ID         string              `json:"-"`
	ActingUser basedto.ObjectIDReq `json:"actingUser"`
}

func NewUpdateAPIKeyReq() *UpdateAPIKeyReq {
	return &UpdateAPIKeyReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAPIKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateObjectIDReq(&req.ActingUser, true, "actingUser")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAPIKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
