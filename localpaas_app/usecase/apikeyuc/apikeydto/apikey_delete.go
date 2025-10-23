package apikeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteAPIKeyReq struct {
	ID string `json:"-"`
}

func NewDeleteAPIKeyReq() *DeleteAPIKeyReq {
	return &DeleteAPIKeyReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteAPIKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteAPIKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
