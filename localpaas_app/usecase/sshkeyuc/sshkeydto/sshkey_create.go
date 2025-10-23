package sshkeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateSSHKeyReq struct {
	*SSHKeyBaseReq
}

type SSHKeyBaseReq struct {
	Name            string                    `json:"name"`
	PrivateKey      string                    `json:"privateKey"`
	ProjectAccesses []*SSHKeyProjectAccessReq `json:"projectAccesses"`
}

type SSHKeyProjectAccessReq struct {
	ID      string `json:"id"`
	Allowed bool   `json:"allowed"`
	// NOTE: this field is used to grant access to a project,
	// but deny access to specific apps within the project
	AppAccesses []*SSHKeyAppAccessReq `json:"appAccesses"`
}

type SSHKeyAppAccessReq struct {
	ID      string `json:"id"`
	Allowed bool   `json:"allowed"`
}

func (req *SSHKeyBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, validateSSHKeyName(&req.Name, true, field+"name")...)
	res = append(res, basedto.ValidateStr(&req.PrivateKey, true, 1, maxKeyLen, "privateKey")...)
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
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
