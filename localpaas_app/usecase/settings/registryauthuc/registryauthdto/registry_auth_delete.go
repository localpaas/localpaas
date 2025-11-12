package registryauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteRegistryAuthReq struct {
	ID string `json:"-"`
}

func NewDeleteRegistryAuthReq() *DeleteRegistryAuthReq {
	return &DeleteRegistryAuthReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteRegistryAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteRegistryAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
