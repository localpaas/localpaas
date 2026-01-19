package gittokendto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type UpdateGitTokenReq struct {
	providers.UpdateSettingReq
	*GitTokenBaseReq
}

func NewUpdateGitTokenReq() *UpdateGitTokenReq {
	return &UpdateGitTokenReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateGitTokenReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateGitTokenResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
