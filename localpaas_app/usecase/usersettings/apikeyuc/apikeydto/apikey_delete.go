package apikeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteAPIKeyReq struct {
	settings.DeleteSettingReq
}

func NewDeleteAPIKeyReq() *DeleteAPIKeyReq {
	return &DeleteAPIKeyReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteAPIKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteAPIKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
