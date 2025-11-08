package userdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type BeginMFATotpSetupReq struct {
	CurrentPasscode string `json:"currentPasscode"`
}

func NewBeginMFATotpSetupReq() *BeginMFATotpSetupReq {
	return &BeginMFATotpSetupReq{}
}

func (req *BeginMFATotpSetupReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStr(&req.CurrentPasscode, false,
		minPasscodeLen, maxPasscodeLen, "currentPasscode")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type BeginMFATotpSetupResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *MFATotpSetupDataResp `json:"data"`
}

type MFATotpSetupDataResp struct {
	Secret    string             `json:"secret"`
	TotpToken string             `json:"totpToken"`
	QRCode    *MFATotpQRCodeResp `json:"qrCode"`
}

type MFATotpQRCodeResp struct {
	DataBase64 string `json:"dataBase64"`
	ImageType  string `json:"imageType"`
	ImageSize  int    `json:"imageSize"`
}
