package secretdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	secretKeyMaxLen   = 1000
	secretValueMaxLen = 10 * 1024 * 1024 // 10MB
)

type CreateSecretReq struct {
	settings.CreateSettingReq
	Key    string `json:"key"`
	Value  string `json:"value"`
	Base64 bool   `json:"base64"`
}

func NewCreateSecretReq() *CreateSecretReq {
	return &CreateSecretReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateSecretReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStr(&req.Key, true, 1, secretKeyMaxLen, "key")...)
	validators = append(validators, basedto.ValidateStr(&req.Value, true, 1, secretValueMaxLen, "value")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateSecretResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
