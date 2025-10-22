package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type GetAppEnvVarsReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
}

func NewGetAppEnvVarsReq() *GetAppEnvVarsReq {
	return &GetAppEnvVarsReq{}
}

func (req *GetAppEnvVarsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppEnvVarsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *AppEnvVarsResp   `json:"data"`
}

type AppEnvVarsResp struct {
	EnvVars [][]string `json:"envVars"`
}

func TransformAppEnvVars(envVars *entity.AppEnvVars) (resp *AppEnvVarsResp, err error) {
	resp = &AppEnvVarsResp{
		EnvVars: [][]string{},
	}
	if envVars != nil && len(envVars.Data) > 0 {
		resp.EnvVars = envVars.Data
	}
	return
}
