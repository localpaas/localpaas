package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateAppSettingsReq struct {
	ProjectID          string                 `json:"-"`
	AppID              string                 `json:"-"`
	EnvVars            EnvVarsReq             `json:"envVars"`
	DeploymentSettings *DeploymentSettingsReq `json:"deploymentSettings"`
	HttpSettings       *HttpSettingsReq       `json:"httpSettings"`
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
	validators = append(validators, req.HttpSettings.validate("httpSettings")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
