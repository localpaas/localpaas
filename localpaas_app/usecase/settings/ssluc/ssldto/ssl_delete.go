package ssldto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteSslReq struct {
	ID string `json:"-"`
}

func NewDeleteSslReq() *DeleteSslReq {
	return &DeleteSslReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteSslResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
