package cronjobdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateCronJobReq struct {
	settings.UpdateSettingReq
	*CronJobBaseReq
}

func NewUpdateCronJobReq() *UpdateCronJobReq {
	return &UpdateCronJobReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateCronJobReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateCronJobResp struct {
	Meta *basedto.Meta `json:"meta"`
}
