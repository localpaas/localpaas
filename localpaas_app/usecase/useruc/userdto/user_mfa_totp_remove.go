package userdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type RemoveMFATotpReq struct {
	Passcode string `json:"passcode"`
}

func NewRemoveMFATotpReq() *RemoveMFATotpReq {
	return &RemoveMFATotpReq{}
}

func (req *RemoveMFATotpReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStr(&req.Passcode, true,
		minPasscodeLen, maxPasscodeLen, "passcode")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type RemoveMFATotpResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
