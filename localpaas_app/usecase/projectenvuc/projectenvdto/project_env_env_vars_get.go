package projectenvdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type GetProjectEnvEnvVarsReq struct {
	ProjectID    string `json:"-"`
	ProjectEnvID string `json:"-"`
}

func NewGetProjectEnvEnvVarsReq() *GetProjectEnvEnvVarsReq {
	return &GetProjectEnvEnvVarsReq{}
}

func (req *GetProjectEnvEnvVarsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectEnvID, true, "projectEnvId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetProjectEnvEnvVarsResp struct {
	Meta *basedto.BaseMeta      `json:"meta"`
	Data *ProjectEnvEnvVarsResp `json:"data"`
}

type ProjectEnvEnvVarsResp struct {
	EnvVars [][]string `json:"envVars"`
}

func TransformProjectEnvEnvVars(envVars *entity.ProjectEnvEnvVars) (resp *ProjectEnvEnvVarsResp, err error) {
	resp = &ProjectEnvEnvVarsResp{
		EnvVars: [][]string{},
	}
	if envVars != nil && len(envVars.Data) > 0 {
		resp.EnvVars = envVars.Data
	}
	return
}
