package userdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	minPasscodeLen = 1
	maxPasscodeLen = 10

	minTokenLen = 10
	maxTokenLen = 10000
)

type CompleteMFATotpSetupReq struct {
	Passcode  string `json:"passcode"`
	TotpToken string `json:"totpToken"`
}

func NewCompleteMFATotpSetupReq() *CompleteMFATotpSetupReq {
	return &CompleteMFATotpSetupReq{}
}

func (req *CompleteMFATotpSetupReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStr(&req.Passcode, true,
		minPasscodeLen, maxPasscodeLen, "passcode")...)
	validators = append(validators, basedto.ValidateStr(&req.TotpToken, true,
		minTokenLen, maxTokenLen, "totpToken")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CompleteMFATotpSetupResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
