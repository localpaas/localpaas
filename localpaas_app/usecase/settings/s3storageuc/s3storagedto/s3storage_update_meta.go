package s3storagedto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateS3StorageMetaReq struct {
	ID       string              `json:"-"`
	Status   *base.SettingStatus `json:"status"`
	ExpireAt *time.Time          `json:"expireAt"`
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
