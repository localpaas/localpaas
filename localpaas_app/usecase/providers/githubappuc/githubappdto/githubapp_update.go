package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateGithubAppReq struct {
	ID        string `json:"-"`
	UpdateVer int    `json:"updateVer"`
	*GithubAppBaseReq
}

func NewUpdateGithubAppReq() *UpdateGithubAppReq {
	return &UpdateGithubAppReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateGithubAppReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateGithubAppResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
