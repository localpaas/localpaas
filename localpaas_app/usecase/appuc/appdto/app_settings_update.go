package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateAppSettingsReq struct {
	ProjectID string          `json:"-"`
	AppID     string          `json:"-"`
	Settings  *AppSettingsReq `json:"settings"`
}

type AppSettingsReq struct {
	Test string `json:"test"`
}

func (p *AppSettingsReq) validate(_ string) []vld.Validator { //nolint
	if p == nil {
		return nil
	}
	// TODO: add validation
	return nil
}

func NewUpdateAppSettingsReq() *UpdateAppSettingsReq {
	return &UpdateAppSettingsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAppSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, req.Settings.validate("settings")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
