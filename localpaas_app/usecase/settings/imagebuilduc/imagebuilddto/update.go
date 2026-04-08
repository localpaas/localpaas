package imagebuilddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateImageBuildReq struct {
	settings.UpdateSettingReq
	*ImageBuildBaseReq
}

func NewUpdateImageBuildReq() *UpdateImageBuildReq {
	return &UpdateImageBuildReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateImageBuildReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateImageBuildResp struct {
	Meta *basedto.Meta `json:"meta"`
}
