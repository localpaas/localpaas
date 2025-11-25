package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

//
// REQUEST
//

type EnvVarsReq struct {
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

func (req *EnvVarsReq) validate(_ string) []vld.Validator { //nolint
	if req == nil {
		return nil
	}
	// TODO: add validation
	return nil
}

//
// RESPONSE
//

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

func TransformEnvVars(input *AppSettingsTransformationInput) (resp *EnvVarsResp, err error) {
	if len(input.EnvVars) == 0 {
		return nil, nil
	}

	var appEnvVars, parentAppEnvVars, projectEnvVars *entity.EnvVars
	for _, env := range input.EnvVars {
		switch env.ObjectID {
		case input.App.ID:
			appEnvVars = env.MustAsEnvVars()
		case input.App.ProjectID:
			projectEnvVars = env.MustAsEnvVars()
		case input.App.ParentID:
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
