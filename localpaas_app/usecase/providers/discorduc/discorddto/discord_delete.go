package discorddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type DeleteDiscordReq struct {
	providers.DeleteSettingReq
}

func NewDeleteDiscordReq() *DeleteDiscordReq {
	return &DeleteDiscordReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteDiscordReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteDiscordResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
