package imagebuilddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteUniqueImageBuildReq struct {
	settings.DeleteUniqueSettingReq
}

func NewDeleteUniqueImageBuildReq() *DeleteUniqueImageBuildReq {
	return &DeleteUniqueImageBuildReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteUniqueImageBuildReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteUniqueSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteUniqueImageBuildResp struct {
	Meta *basedto.Meta `json:"meta"`
}
