package s3storagedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateS3StorageReq struct {
	settings.UpdateSettingReq
	*S3StorageBaseReq
}

func NewUpdateS3StorageReq() *UpdateS3StorageReq {
	return &UpdateS3StorageReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateS3StorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateS3StorageResp struct {
	Meta *basedto.Meta `json:"meta"`
}
