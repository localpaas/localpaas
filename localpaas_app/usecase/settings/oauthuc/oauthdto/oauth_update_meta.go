package oauthdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateOAuthMetaReq struct {
	ID       string              `json:"-"`
	Status   *base.SettingStatus `json:"status"`
	ExpireAt *time.Time          `json:"expireAt"`
}

func NewUpdateOAuthMetaReq() *UpdateOAuthMetaReq {
	return &UpdateOAuthMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateOAuthMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false, base.AllSettingStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateOAuthMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
