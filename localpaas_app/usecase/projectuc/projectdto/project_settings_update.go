package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type UpdateProjectSettingsReq struct {
	ProjectID string              `json:"-"`
	Settings  *GeneralSettingsReq `json:"settings"`
	EnvVars   EnvVarsReq          `json:"envVars"`
}

type EnvVarsReq []*EnvVarReq

type EnvVarReq struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsBuildEnv bool   `json:"isBuildEnv"`
}

func (req *EnvVarReq) ToEntity() *entity.EnvVar {
	return &entity.EnvVar{
		Key:        req.Key,
		Value:      req.Value,
		IsBuildEnv: req.IsBuildEnv,
	}
}

func (req *EnvVarsReq) validate(_ string) []vld.Validator { //nolint
	if req == nil {
		return nil
	}
	// TODO: add validation
	return nil
}

type GeneralSettingsReq struct {
	Test string `json:"test"`
}

func (req *GeneralSettingsReq) ToEntity() *entity.ProjectSettings {
	return &entity.ProjectSettings{
		Test: req.Test,
	}
}

func (req *GeneralSettingsReq) validate(_ string) []vld.Validator { //nolint
	if req == nil {
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
	validators = append(validators, req.EnvVars.validate("envVars")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateProjectSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
