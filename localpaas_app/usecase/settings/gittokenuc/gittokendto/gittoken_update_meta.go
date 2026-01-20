package gittokendto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateGitTokenMetaReq struct {
	settings.UpdateSettingMetaReq
}

func NewUpdateGitTokenMetaReq() *UpdateGitTokenMetaReq {
	return &UpdateGitTokenMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateGitTokenMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateGitTokenMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
