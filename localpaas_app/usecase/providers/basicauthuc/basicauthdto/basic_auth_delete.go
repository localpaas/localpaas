package basicauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type DeleteBasicAuthReq struct {
	providers.DeleteSettingReq
}

func NewDeleteBasicAuthReq() *DeleteBasicAuthReq {
	return &DeleteBasicAuthReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteBasicAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteBasicAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
