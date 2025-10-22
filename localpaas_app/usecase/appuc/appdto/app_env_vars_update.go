package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateAppEnvVarsReq struct {
	ProjectID string     `json:"-"`
	AppID     string     `json:"-"`
	EnvVars   [][]string `json:"envVars"`
}

func NewUpdateAppEnvVarsReq() *UpdateAppEnvVarsReq {
	return &UpdateAppEnvVarsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAppEnvVarsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	// TODO: add validation for req.EnvVars
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppEnvVarsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
