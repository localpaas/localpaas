package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateProjectSettingsReq struct {
	ProjectID string              `json:"-"`
	Settings  *ProjectSettingsReq `json:"settings"`
}

type ProjectSettingsReq struct {
	Test string `json:"test"`
}

func (p *ProjectSettingsReq) validate(_ string) []vld.Validator { //nolint
	if p == nil {
		return nil
	}
	// TODO: add validation
	return nil
}

func NewUpdateProjectSettingsReq() *UpdateProjectSettingsReq {
	return &UpdateProjectSettingsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateProjectSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, req.Settings.validate("settings")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateProjectSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
