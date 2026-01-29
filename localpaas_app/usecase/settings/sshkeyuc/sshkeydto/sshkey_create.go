package sshkeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateSSHKeyReq struct {
	settings.CreateSettingReq
	*SSHKeyBaseReq
}

type SSHKeyBaseReq struct {
	Name       string `json:"name"`
	PrivateKey string `json:"privateKey"`
	Passphrase string `json:"passphrase"`
}

func (req *SSHKeyBaseReq) ToEntity() *entity.SSHKey {
	return &entity.SSHKey{
		PrivateKey: entity.NewEncryptedField(req.PrivateKey),
		Passphrase: entity.NewEncryptedField(req.Passphrase),
	}
}

func (req *SSHKeyBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, validateSSHKeyName(&req.Name, true, field+"name")...)
	res = append(res, basedto.ValidateStr(&req.PrivateKey, true, 1, maxKeyLen, "privateKey")...)
	res = append(res, basedto.ValidateStr(&req.Passphrase, false, 1, maxNameLen, "passphrase")...)
	return res
}

func NewCreateSSHKeyReq() *CreateSSHKeyReq {
	return &CreateSSHKeyReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateSSHKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateSSHKeyResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
