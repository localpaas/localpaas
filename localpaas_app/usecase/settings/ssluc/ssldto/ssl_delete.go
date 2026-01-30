package ssldto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteSSLReq struct {
	settings.DeleteSettingReq
}

func NewDeleteSSLReq() *DeleteSSLReq {
	return &DeleteSSLReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteSSLReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteSSLResp struct {
	Meta *basedto.Meta `json:"meta"`
}
