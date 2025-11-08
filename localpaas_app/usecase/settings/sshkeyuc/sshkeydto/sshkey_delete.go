package sshkeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteSSHKeyReq struct {
	ID string `json:"-"`
}

func NewDeleteSSHKeyReq() *DeleteSSHKeyReq {
	return &DeleteSSHKeyReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteSSHKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteSSHKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
