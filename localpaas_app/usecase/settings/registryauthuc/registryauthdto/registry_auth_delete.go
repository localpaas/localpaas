package registryauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteRegistryAuthReq struct {
	settings.DeleteSettingReq
}

func NewDeleteRegistryAuthReq() *DeleteRegistryAuthReq {
	return &DeleteRegistryAuthReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteRegistryAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteRegistryAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
