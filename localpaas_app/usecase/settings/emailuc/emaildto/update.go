package emaildto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateEmailReq struct {
	settings.UpdateSettingReq
	*EmailBaseReq
}

func NewUpdateEmailReq() *UpdateEmailReq {
	return &UpdateEmailReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateEmailReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateEmailResp struct {
	Meta *basedto.Meta `json:"meta"`
}
