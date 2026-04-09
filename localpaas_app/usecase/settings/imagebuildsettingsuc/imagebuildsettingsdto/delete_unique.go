package imagebuildsettingsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteUniqueImageBuildSettingsReq struct {
	settings.DeleteUniqueSettingReq
}

func NewDeleteUniqueImageBuildSettingsReq() *DeleteUniqueImageBuildSettingsReq {
	return &DeleteUniqueImageBuildSettingsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteUniqueImageBuildSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteUniqueSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteUniqueImageBuildSettingsResp struct {
	Meta *basedto.Meta `json:"meta"`
}
