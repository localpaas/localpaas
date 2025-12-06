package gittokendto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateGitTokenReq struct {
	ID        string `json:"-"`
	UpdateVer int    `json:"updateVer"`
	*GitTokenBaseReq
}

func NewUpdateGitTokenReq() *UpdateGitTokenReq {
	return &UpdateGitTokenReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateGitTokenReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateGitTokenResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
