package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteGithubAppReq struct {
	ID string `json:"-"`
}

func NewDeleteGithubAppReq() *DeleteGithubAppReq {
	return &DeleteGithubAppReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteGithubAppReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteGithubAppResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
