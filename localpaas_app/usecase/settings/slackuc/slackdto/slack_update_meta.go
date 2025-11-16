package slackdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateSlackMetaReq struct {
	ID       string              `json:"-"`
	Status   *base.SettingStatus `json:"status"`
	ExpireAt *time.Time          `json:"expireAt"`
}

func NewUpdateSlackMetaReq() *UpdateSlackMetaReq {
	return &UpdateSlackMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateSlackMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false, base.AllSettingStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateSlackMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
