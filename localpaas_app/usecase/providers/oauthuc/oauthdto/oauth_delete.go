package oauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type DeleteOAuthReq struct {
	providers.DeleteSettingReq
}

func NewDeleteOAuthReq() *DeleteOAuthReq {
	return &DeleteOAuthReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteOAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteOAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
