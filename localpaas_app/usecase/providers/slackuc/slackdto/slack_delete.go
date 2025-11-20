package slackdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteSlackReq struct {
	ID string `json:"-"`
}

func NewDeleteSlackReq() *DeleteSlackReq {
	return &DeleteSlackReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteSlackReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteSlackResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
