package s3storagedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type DeleteS3StorageReq struct {
	providers.DeleteSettingReq
}

func NewDeleteS3StorageReq() *DeleteS3StorageReq {
	return &DeleteS3StorageReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteS3StorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteS3StorageResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
