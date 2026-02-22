package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type SetupGithubAppManifestFlowReq struct {
	SettingID      string `json:"-" mapstructure:"-"`
	Code           string `json:"-" mapstructure:"code"`
	State          string `json:"-" mapstructure:"state"`
	InstallationID int64  `json:"-" mapstructure:"installation_id"`
	SetupAction    string `json:"-" mapstructure:"setup_action"`
}

func NewSetupGithubAppManifestFlowReq() *SetupGithubAppManifestFlowReq {
	return &SetupGithubAppManifestFlowReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *SetupGithubAppManifestFlowReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type SetupGithubAppManifestFlowResp struct {
	Meta *basedto.Meta                       `json:"meta"`
	Data *SetupGithubAppManifestFlowDataResp `json:"data"`
}

type SetupGithubAppManifestFlowDataResp struct {
	RedirectURL string `json:"redirectURL"`
}
