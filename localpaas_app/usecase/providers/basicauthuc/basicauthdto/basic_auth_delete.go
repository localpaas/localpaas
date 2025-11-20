package basicauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteBasicAuthReq struct {
	ID string `json:"-"`
}

func NewDeleteBasicAuthReq() *DeleteBasicAuthReq {
	return &DeleteBasicAuthReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteBasicAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteBasicAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
