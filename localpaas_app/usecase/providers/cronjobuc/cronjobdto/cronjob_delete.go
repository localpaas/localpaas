package cronjobdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type DeleteCronJobReq struct {
	providers.DeleteSettingReq
}

func NewDeleteCronJobReq() *DeleteCronJobReq {
	return &DeleteCronJobReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteCronJobReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteCronJobResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
