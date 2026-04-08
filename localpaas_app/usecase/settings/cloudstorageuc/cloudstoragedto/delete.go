package cloudstoragedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteCloudStorageReq struct {
	settings.DeleteSettingReq
}

func NewDeleteCloudStorageReq() *DeleteCloudStorageReq {
	return &DeleteCloudStorageReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteCloudStorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteCloudStorageResp struct {
	Meta *basedto.Meta `json:"meta"`
}
