package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateProjectEnvVarsReq struct {
	ProjectID string     `json:"-"`
	EnvVars   [][]string `json:"envVars"`
}

func NewUpdateProjectEnvVarsReq() *UpdateProjectEnvVarsReq {
	return &UpdateProjectEnvVarsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateProjectEnvVarsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	// TODO: add validation for req.EnvVars
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateProjectEnvVarsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
