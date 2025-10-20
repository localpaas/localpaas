package sessiondto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type LoginWithPasscodeReq struct {
	Passcode string `json:"passcode"`
	MFAToken string `json:"mfaToken"`
}

func NewLoginWithPasscodeReq() *LoginWithPasscodeReq {
	return &LoginWithPasscodeReq{}
}

func (req *LoginWithPasscodeReq) Validate() apperrors.ValidationErrors {
	return apperrors.NewValidationErrors(vld.Validate(
		vld.Required(req.Passcode).OnError(
			vld.SetField("passcode", nil),
			vld.SetCustomKey("ERR_VLD_FIELD_REQUIRED"),
		),
		vld.Required(req.MFAToken).OnError(
			vld.SetField("mfaToken", nil),
			vld.SetCustomKey("ERR_VLD_FIELD_REQUIRED"),
		),
	))
}

type LoginWithPasscodeResp struct {
	Meta *basedto.BaseMeta          `json:"meta"`
	Data *LoginWithPasscodeDataResp `json:"data"`
}

type LoginWithPasscodeDataResp struct {
	Session *BaseCreateSessionResp `json:"session,omitempty"`
}
