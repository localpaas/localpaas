package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetAppSettingsReq struct {
	ProjectID string             `json:"-"`
	AppID     string             `json:"-"`
	Type      []base.SettingType `json:"-" mapstructure:"type"`
}

func NewGetAppSettingsReq() *GetAppSettingsReq {
	return &GetAppSettingsReq{}
}

func (req *GetAppSettingsReq) Validate() apperrors.ValidationErrors {
	if len(req.Type) == 0 {
		req.Type = base.AllAppSettingTypes
	}

	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateSlice(req.Type, true, 0, base.AllAppSettingTypes, "type")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *AppSettingsResp  `json:"data"`
}

type AppSettingsResp struct {
	EnvVars            *EnvVarsResp            `json:"envVars,omitempty"`
	DeploymentSettings *DeploymentSettingsResp `json:"deploymentSettings,omitempty"`
}

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

type DeploymentSettingsResp struct {
	Test string `json:"test"`
}

func TransformAppSettings(app *entity.App) (resp *AppSettingsResp, err error) {
	resp = &AppSettingsResp{}

	var envVarSettings []*entity.Setting
	for _, setting := range app.Settings {
		switch setting.Type { //nolint:exhaustive
		case base.SettingTypeEnvVar:
			envVarSettings = append(envVarSettings, setting)
		case base.SettingTypeDeployment:
			resp.DeploymentSettings, err = TransformDeploymentSettings(setting)
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
		}
	}

	if len(envVarSettings) > 0 {
		resp.EnvVars, err = TransformEnvVars(app, envVarSettings)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	return resp, nil
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

func TransformDeploymentSettings(setting *entity.Setting) (resp *DeploymentSettingsResp, err error) {
	data, err := setting.ParseAppDeploymentSettings()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, &data); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
