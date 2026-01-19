package discorddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type UpdateDiscordReq struct {
	providers.UpdateSettingReq
	*DiscordBaseReq
}

func NewUpdateDiscordReq() *UpdateDiscordReq {
	return &UpdateDiscordReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateDiscordReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateDiscordResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
