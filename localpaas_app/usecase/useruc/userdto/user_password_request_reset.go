package userdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type RequestResetPasswordReq struct {
	ID                 string `json:"-"`
	SendResettingEmail bool   `json:"sendResettingEmail"`
}

func NewRequestResetPasswordReq() *RequestResetPasswordReq {
	return &RequestResetPasswordReq{}
}

func (req *RequestResetPasswordReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type RequestResetPasswordResp struct {
	Meta *basedto.BaseMeta             `json:"meta"`
	Data *RequestResetPasswordDataResp `json:"data"`
}

type RequestResetPasswordDataResp struct {
	ResetPasswordLink string `json:"resetPasswordLink"`
}
