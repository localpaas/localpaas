package gittokendto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteGitTokenReq struct {
	ID string `json:"-"`
}

func NewDeleteGitTokenReq() *DeleteGitTokenReq {
	return &DeleteGitTokenReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteGitTokenReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteGitTokenResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
