package oauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteOAuthReq struct {
	ID string `json:"-"`
}

func NewDeleteOAuthReq() *DeleteOAuthReq {
	return &DeleteOAuthReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteOAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteOAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
