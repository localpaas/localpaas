package sessiondto

import (
	"strings"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/translation"
)

type LoginWithPasswordReq struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	TrustedDeviceID string `json:"trustedDeviceId"`

	AcceptLanguage translation.Lang `json:"-"`
}

func NewLoginWithPasswordReq() *LoginWithPasswordReq {
	return &LoginWithPasswordReq{}
}

func (req *LoginWithPasswordReq) Validate() apperrors.ValidationErrors {
	req.Username = strings.TrimSpace(req.Username)
	req.TrustedDeviceID = strings.TrimSpace(req.TrustedDeviceID)

	return apperrors.NewValidationErrors(vld.Validate(
		vld.Required(req.Username).OnError(
			vld.SetField("username", nil),
			vld.SetCustomKey("ERR_VLD_FIELD_REQUIRED"),
		),
		vld.Required(req.Password).OnError(
			vld.SetField("password", nil),
			vld.SetCustomKey("ERR_VLD_FIELD_REQUIRED"),
		),
	))
}

type LoginWithPasswordResp struct {
	Meta *basedto.BaseMeta          `json:"meta"`
	Data *LoginWithPasswordDataResp `json:"data"`
}

type LoginWithPasswordDataResp struct {
	NextStep string                 `json:"nextStep,omitempty"`
	MFAType  base.MFAType           `json:"mfaType,omitempty"`
	MFAToken string                 `json:"mfaToken,omitempty"`
	Session  *BaseCreateSessionResp `json:"session,omitempty"`
}
