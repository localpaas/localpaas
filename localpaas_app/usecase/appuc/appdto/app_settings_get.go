package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
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
	HttpSettings       *HttpSettingsResp       `json:"httpSettings,omitempty"`
}

func TransformAppSettings(app *entity.App) (resp *AppSettingsResp, err error) {
	resp = &AppSettingsResp{}

	var envVarSettings []*entity.Setting
	for _, setting := range app.Settings {
		switch setting.Type { //nolint:exhaustive
		case base.SettingTypeEnvVar:
			envVarSettings = append(envVarSettings, setting)

		case base.SettingTypeAppDeployment:
			resp.DeploymentSettings, err = TransformDeploymentSettings(setting)
			if err != nil {
				return nil, apperrors.Wrap(err)
			}

		case base.SettingTypeAppHttp:
			resp.HttpSettings, err = TransformHttpSettings(setting)
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
