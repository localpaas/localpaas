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

type AppSettingsTransformationInput struct {
	App                *entity.App
	EnvVars            []*entity.EnvVars
	DeploymentSettings *entity.AppDeploymentSettings
	HttpSettings       *entity.AppHttpSettings

	DefaultNginxSettings *entity.NginxSettings
	ReferenceSettingMap  map[string]*entity.Setting
}

func TransformAppSettings(input *AppSettingsTransformationInput) (resp *AppSettingsResp, err error) {
	resp = &AppSettingsResp{}

	resp.DeploymentSettings, err = TransformDeploymentSettings(input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.HttpSettings, err = TransformHttpSettings(input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.EnvVars, err = TransformEnvVars(input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}
