package discorddto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateDiscordMetaReq struct {
	ID        string              `json:"-"`
	Status    *base.SettingStatus `json:"status"`
	ExpireAt  *time.Time          `json:"expireAt"`
	UpdateVer int                 `json:"updateVer"`
}

func NewUpdateDiscordMetaReq() *UpdateDiscordMetaReq {
	return &UpdateDiscordMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateDiscordMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateDiscordMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
