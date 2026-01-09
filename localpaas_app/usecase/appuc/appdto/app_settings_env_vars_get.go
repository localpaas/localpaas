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
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppEnvVarsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *EnvVarsResp      `json:"data"`
}

type EnvVarsResp struct {
	App       []*EnvVarResp `json:"app"`
	ParentApp []*EnvVarResp `json:"parentApp"`
	Project   []*EnvVarResp `json:"project"`
	UpdateVer int           `json:"updateVer"`
}

type EnvVarResp struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsBuildEnv bool   `json:"isBuildEnv,omitempty"`
}

func TransformEnvVars(app *entity.App, envVars []*entity.Setting) (resp *EnvVarsResp, err error) {
	if len(envVars) == 0 {
		return nil, nil
	}

	var appEnvVars, parentAppEnvVars, projectEnvVars *entity.EnvVars
	for _, env := range envVars {
		switch env.ObjectID {
		case app.ID:
			appEnvVars = env.MustAsEnvVars()
		case app.ProjectID:
			projectEnvVars = env.MustAsEnvVars()
		case app.ParentID:
			parentAppEnvVars = env.MustAsEnvVars()
		}
	}

	resp = &EnvVarsResp{
		App:       []*EnvVarResp{},
		ParentApp: []*EnvVarResp{},
		Project:   []*EnvVarResp{},
	}
	if appEnvVars != nil {
		for _, v := range appEnvVars.Data {
			resp.App = append(resp.App, &EnvVarResp{
				Key:        v.Key,
				Value:      v.Value,
				IsBuildEnv: v.IsBuildEnv,
			})
		}
	}
	if parentAppEnvVars != nil {
		for _, v := range parentAppEnvVars.Data {
			resp.ParentApp = append(resp.ParentApp, &EnvVarResp{
				Key:        v.Key,
				Value:      v.Value,
				IsBuildEnv: v.IsBuildEnv,
			})
		}
	}
	if projectEnvVars != nil {
		for _, v := range projectEnvVars.Data {
			resp.Project = append(resp.Project, &EnvVarResp{
				Key:        v.Key,
				Value:      v.Value,
				IsBuildEnv: v.IsBuildEnv,
			})
		}
	}
	return resp, nil
}
