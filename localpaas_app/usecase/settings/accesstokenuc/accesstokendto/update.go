package accesstokendto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateAccessTokenReq struct {
	settings.UpdateSettingReq
	*AccessTokenBaseReq
}

func NewUpdateAccessTokenReq() *UpdateAccessTokenReq {
	return &UpdateAccessTokenReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAccessTokenReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAccessTokenResp struct {
	Meta *basedto.Meta `json:"meta"`
}
