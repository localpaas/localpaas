package cloudproviderdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteCloudProviderReq struct {
	settings.DeleteSettingReq
}

func NewDeleteCloudProviderReq() *DeleteCloudProviderReq {
	return &DeleteCloudProviderReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteCloudProviderReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteCloudProviderResp struct {
	Meta *basedto.Meta `json:"meta"`
}
