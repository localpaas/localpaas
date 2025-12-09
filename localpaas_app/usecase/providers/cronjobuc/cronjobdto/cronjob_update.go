package cronjobdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateCronJobReq struct {
	ID        string `json:"-"`
	UpdateVer int    `json:"updateVer"`
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
	Meta *basedto.BaseMeta `json:"meta"`
}
