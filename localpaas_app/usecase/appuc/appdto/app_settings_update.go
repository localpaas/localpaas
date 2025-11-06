package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type UpdateAppSettingsReq struct {
	ProjectID          string                 `json:"-"`
	AppID              string                 `json:"-"`
	EnvVars            EnvVarsReq             `json:"envVars"`
	DeploymentSettings *DeploymentSettingsReq `json:"deploymentSettings"`
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

type DeploymentSettingsReq struct {
	Test string `json:"test"`
}

func (req *DeploymentSettingsReq) ToEntity() *entity.AppDeploymentSettings {
	return &entity.AppDeploymentSettings{
		Test: req.Test,
	}
}

func (req *DeploymentSettingsReq) validate(_ string) []vld.Validator { //nolint
	if req == nil {
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
	validators = append(validators, req.EnvVars.validate("envVars")...)
	validators = append(validators, req.DeploymentSettings.validate("deploymentSettings")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
