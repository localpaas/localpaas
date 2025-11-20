package secretdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteSecretReq struct {
	ID       string `json:"-"`
	ObjectID string `json:"-"`
}

func NewDeleteSecretReq() *DeleteSecretReq {
	return &DeleteSecretReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteSecretReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteSecretResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
