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
	*SecretBaseReq
}

type SecretBaseReq struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Base64 bool   `json:"base64"`
}

func (req *SecretBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Key, true, 1, secretKeyMaxLen, field+"key")...)
	res = append(res, basedto.ValidateStr(&req.Value, true, 1, secretValueMaxLen, field+"value")...)
	return res
}

func NewCreateSecretReq() *CreateSecretReq {
	return &CreateSecretReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateSecretReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateSecretResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
