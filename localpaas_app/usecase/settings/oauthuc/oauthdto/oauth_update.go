package oauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateOAuthReq struct {
	settings.UpdateSettingReq
	*OAuthBaseReq
}

func NewUpdateOAuthReq() *UpdateOAuthReq {
	return &UpdateOAuthReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateOAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateOAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
