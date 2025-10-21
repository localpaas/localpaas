package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type GetProjectEnvVarsReq struct {
	ProjectID string `json:"-"`
}

func NewGetProjectEnvVarsReq() *GetProjectEnvVarsReq {
	return &GetProjectEnvVarsReq{}
}

func (req *GetProjectEnvVarsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetProjectEnvVarsResp struct {
	Meta *basedto.BaseMeta   `json:"meta"`
	Data *ProjectEnvVarsResp `json:"data"`
}

type ProjectEnvVarsResp struct {
	EnvVars [][]string `json:"envVars"`
}

func TransformProjectEnvVars(envVars *entity.ProjectEnvVars) (resp *ProjectEnvVarsResp, err error) {
	resp = &ProjectEnvVarsResp{
		EnvVars: [][]string{},
	}
	if envVars != nil && len(envVars.Data) > 0 {
		resp.EnvVars = envVars.Data
	}
	return
}
