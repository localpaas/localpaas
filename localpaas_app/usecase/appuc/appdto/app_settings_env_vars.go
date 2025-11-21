package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

//
// REQUEST
//

type EnvVarsReq []*EnvVarReq

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
}

type EnvVarResp struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsBuildEnv bool   `json:"isBuildEnv,omitempty"`
}

func TransformEnvVars(app *entity.App, envSettings []*entity.Setting) (resp *EnvVarsResp, err error) {
	var appEnvVars, parentAppEnvVars, projectEnvVars *entity.EnvVars
	for _, envSetting := range envSettings {
		switch envSetting.ObjectID {
		case app.ID:
			appEnvVars, err = envSetting.ParseEnvVars()
		case app.ProjectID:
			projectEnvVars, err = envSetting.ParseEnvVars()
		case app.ParentID:
			parentAppEnvVars, err = envSetting.ParseEnvVars()
		}
		if err != nil {
			return nil, apperrors.Wrap(err)
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
