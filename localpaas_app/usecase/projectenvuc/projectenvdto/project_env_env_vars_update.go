package projectenvdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateProjectEnvEnvVarsReq struct {
	ProjectID    string     `json:"-"`
	ProjectEnvID string     `json:"-"`
	EnvVars      [][]string `json:"envVars"`
}

func NewUpdateProjectEnvEnvVarsReq() *UpdateProjectEnvEnvVarsReq {
	return &UpdateProjectEnvEnvVarsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateProjectEnvEnvVarsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectEnvID, true, "projectEnvId")...)
	// TODO: add validation for req.EnvVars
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateProjectEnvEnvVarsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
