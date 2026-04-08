package imagebuilddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteImageBuildReq struct {
	settings.DeleteSettingReq
}

func NewDeleteImageBuildReq() *DeleteImageBuildReq {
	return &DeleteImageBuildReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteImageBuildReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteImageBuildResp struct {
	Meta *basedto.Meta `json:"meta"`
}
