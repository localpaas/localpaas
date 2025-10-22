package projectenvdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateProjectEnvSettingsReq struct {
	ProjectID    string                 `json:"-"`
	ProjectEnvID string                 `json:"-"`
	Settings     *ProjectEnvSettingsReq `json:"settings"`
}

type ProjectEnvSettingsReq struct {
	Test string `json:"test"`
}

func (p *ProjectEnvSettingsReq) validate(_ string) []vld.Validator { //nolint
	if p == nil {
		return nil
	}
	// TODO: add validation
	return nil
}

func NewUpdateProjectEnvSettingsReq() *UpdateProjectEnvSettingsReq {
	return &UpdateProjectEnvSettingsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateProjectEnvSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectEnvID, true, "projectEnvId")...)
	validators = append(validators, req.Settings.validate("settings")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateProjectEnvSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
