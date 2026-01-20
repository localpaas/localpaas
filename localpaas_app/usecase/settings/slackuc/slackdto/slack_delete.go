package slackdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteSlackReq struct {
	settings.DeleteSettingReq
}

func NewDeleteSlackReq() *DeleteSlackReq {
	return &DeleteSlackReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteSlackReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteSlackResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
