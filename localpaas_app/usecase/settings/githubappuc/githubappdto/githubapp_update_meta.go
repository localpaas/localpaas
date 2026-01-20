package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateGithubAppMetaReq struct {
	settings.UpdateSettingMetaReq
}

func NewUpdateGithubAppMetaReq() *UpdateGithubAppMetaReq {
	return &UpdateGithubAppMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateGithubAppMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllSettingSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateGithubAppMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
