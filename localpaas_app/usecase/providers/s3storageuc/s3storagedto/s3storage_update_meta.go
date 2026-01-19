package s3storagedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type UpdateS3StorageMetaReq struct {
	providers.UpdateSettingMetaReq
}

func NewUpdateS3StorageMetaReq() *UpdateS3StorageMetaReq {
	return &UpdateS3StorageMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateS3StorageMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateS3StorageMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
