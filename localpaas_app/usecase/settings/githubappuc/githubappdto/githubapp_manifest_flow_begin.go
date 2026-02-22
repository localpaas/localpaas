package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type BeginGithubAppManifestFlowReq struct {
	settings.CreateSettingReq
	Name string `json:"name"`
	Org  string `json:"org"`
}

func NewBeginGithubAppManifestFlowReq() *BeginGithubAppManifestFlowReq {
	return &BeginGithubAppManifestFlowReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *BeginGithubAppManifestFlowReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type BeginGithubAppManifestFlowResp struct {
	Meta *basedto.Meta                       `json:"meta"`
	Data *BeginGithubAppManifestFlowDataResp `json:"data"`
}

type BeginGithubAppManifestFlowDataResp struct {
	RedirectURL string `json:"redirectURL"`
}
