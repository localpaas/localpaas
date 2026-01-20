package slackdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateSlackReq struct {
	settings.UpdateSettingReq
	*SlackBaseReq
}

func NewUpdateSlackReq() *UpdateSlackReq {
	return &UpdateSlackReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateSlackReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateSlackResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
