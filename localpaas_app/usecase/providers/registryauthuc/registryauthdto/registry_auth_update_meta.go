package registryauthdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateRegistryAuthMetaReq struct {
	ID        string              `json:"-"`
	Status    *base.SettingStatus `json:"status"`
	ExpireAt  *time.Time          `json:"expireAt"`
	UpdateVer int                 `json:"updateVer"`
}

func NewUpdateRegistryAuthMetaReq() *UpdateRegistryAuthMetaReq {
	return &UpdateRegistryAuthMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateRegistryAuthMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateRegistryAuthMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
