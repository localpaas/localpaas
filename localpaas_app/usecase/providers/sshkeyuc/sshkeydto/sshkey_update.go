package sshkeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateSSHKeyReq struct {
	ID        string `json:"-"`
	UpdateVer int    `json:"updateVer"`
	*SSHKeyPartialReq
}

type SSHKeyPartialReq struct {
	Name            *string                   `json:"name"`
	PrivateKey      *string                   `json:"privateKey"`
	Passphrase      *string                   `json:"passphrase"`
	ProjectAccesses []*SSHKeyProjectAccessReq `json:"projectAccesses"`
}

func (req *SSHKeyPartialReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, validateSSHKeyName(req.Name, false, field+"name")...)
	res = append(res, basedto.ValidateStr(req.PrivateKey, false, 1, maxKeyLen, "privateKey")...)
	res = append(res, basedto.ValidateStr(req.Passphrase, false, 1, maxNameLen, "passphrase")...)
	return res
}

func NewUpdateSSHKeyReq() *UpdateSSHKeyReq {
	return &UpdateSSHKeyReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateSSHKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateSSHKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
