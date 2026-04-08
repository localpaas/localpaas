package cronjobdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ExecuteCronJobReq struct {
	settings.GetSettingReq
}

func NewExecuteCronJobReq() *ExecuteCronJobReq {
	return &ExecuteCronJobReq{}
}

func (req *ExecuteCronJobReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ExecuteCronJobResp struct {
	Meta *basedto.Meta           `json:"meta"`
	Data *ExecuteCronJobDataResp `json:"data"`
}

type ExecuteCronJobDataResp struct {
	Task *basedto.ObjectIDResp `json:"task"`
}
