package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/secretuc/secretdto"
)

type GetProjectSettingsReq struct {
	ProjectID string             `json:"-"`
	Type      []base.SettingType `json:"-" mapstructure:"type"`
}

func NewGetProjectSettingsReq() *GetProjectSettingsReq {
	return &GetProjectSettingsReq{}
}

func (req *GetProjectSettingsReq) Validate() apperrors.ValidationErrors {
	if len(req.Type) == 0 {
		req.Type = base.AllProjectSettingTypes
	}

	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateSlice(req.Type, true, 0, base.AllProjectSettingTypes, "type")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetProjectSettingsResp struct {
	Meta *basedto.BaseMeta    `json:"meta"`
	Data *ProjectSettingsResp `json:"data"`
}

type ProjectSettingsResp struct {
	EnvVars  *EnvVarsResp         `json:"envVars,omitempty"`
	Secrets  *SecretsResp         `json:"secrets,omitempty"`
	Settings *GeneralSettingsResp `json:"settings,omitempty"`
}

type EnvVarsResp struct {
	Project []*EnvVarResp `json:"project"`
}

type EnvVarResp struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsBuildEnv bool   `json:"isBuildEnv,omitempty"`
}

type SecretsResp struct {
	Project []*secretdto.SecretResp `json:"project"`
}

type GeneralSettingsResp struct {
	Test string `json:"test"`
}

func TransformProjectSettings(project *entity.Project) (resp *ProjectSettingsResp, err error) {
	resp = &ProjectSettingsResp{}

	var allSecrets []*entity.Setting
	for _, setting := range project.Settings {
		switch setting.Type { //nolint:exhaustive
		case base.SettingTypeProject:
			resp.Settings, err = TransformGeneralSettings(setting)
			if err != nil {
				return nil, apperrors.Wrap(err)
			}

		case base.SettingTypeEnvVar:
			resp.EnvVars, err = TransformEnvVars(setting)
			if err != nil {
				return nil, apperrors.Wrap(err)
			}

		case base.SettingTypeSecret:
			allSecrets = append(allSecrets, setting)
		}
	}

	if len(allSecrets) > 0 {
		secrets, err := secretdto.TransformSecrets(allSecrets)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp.Secrets = &SecretsResp{
			Project: secrets,
		}
	}

	return resp, nil
}

func TransformEnvVars(setting *entity.Setting) (resp *EnvVarsResp, err error) {
	if setting == nil {
		return nil, nil
	}
	envVars := setting.MustAsEnvVars()
	if envVars != nil {
		resp = &EnvVarsResp{
			Project: []*EnvVarResp{},
		}
		for _, v := range envVars.Data {
			resp.Project = append(resp.Project, &EnvVarResp{
				Key:        v.Key,
				Value:      v.Value,
				IsBuildEnv: v.IsBuildEnv,
			})
		}
	}
	return resp, nil
}

func TransformGeneralSettings(setting *entity.Setting) (resp *GeneralSettingsResp, err error) {
	data, err := setting.AsProjectSettings()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, &data); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
