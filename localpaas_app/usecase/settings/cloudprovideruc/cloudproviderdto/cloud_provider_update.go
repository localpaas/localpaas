package cloudproviderdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateCloudProviderReq struct {
	settings.UpdateSettingReq
	*CloudProviderBaseReq
}

func NewUpdateCloudProviderReq() *UpdateCloudProviderReq {
	return &UpdateCloudProviderReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateCloudProviderReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateCloudProviderResp struct {
	Meta *basedto.Meta `json:"meta"`
}
