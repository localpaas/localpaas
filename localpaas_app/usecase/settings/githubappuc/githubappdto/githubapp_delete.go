package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteGithubAppReq struct {
	settings.DeleteSettingReq
}

func NewDeleteGithubAppReq() *DeleteGithubAppReq {
	return &DeleteGithubAppReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteGithubAppReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteGithubAppResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
