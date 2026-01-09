package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type UpdateAppEnvVarsReq struct {
	ProjectID string       `json:"-"`
	AppID     string       `json:"-"`
	EnvVars   []*EnvVarReq `json:"envVars"`
	UpdateVer int          `json:"updateVer"`
}

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

func NewUpdateAppEnvVarsReq() *UpdateAppEnvVarsReq {
	return &UpdateAppEnvVarsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAppEnvVarsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	// TODO: validate env var input
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppEnvVarsResp struct {
	Meta *basedto.BaseMeta         `json:"meta"`
	Data *UpdateAppEnvVarsDataResp `json:"data"`
}

type UpdateAppEnvVarsDataResp struct {
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
