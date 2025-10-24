package userdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type BeginMFATotpSetupReq struct {
}

func NewBeginMFATotpSetupReq() *BeginMFATotpSetupReq {
	return &BeginMFATotpSetupReq{}
}

func (req *BeginMFATotpSetupReq) Validate() apperrors.ValidationErrors {
	return nil
}

type BeginMFATotpSetupResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *MFATotpSetupDataResp `json:"data"`
}

type MFATotpSetupDataResp struct {
	Secret    string                  `json:"secret"`
	TotpToken string                  `json:"totpToken"`
	QRCode    *MFATotpSetupQRCodeResp `json:"qrCode"`
}

type MFATotpSetupQRCodeResp struct {
	DataBase64 string `json:"dataBase64"`
	ImageType  string `json:"imageType"`
	ImageSize  int    `json:"imageSize"`
}
