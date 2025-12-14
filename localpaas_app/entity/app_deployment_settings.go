package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentAppDeploymentSettingsVersion = 1
)

type AppDeploymentSettings struct {
	ImageSource *DeploymentImageSource `json:"imageSource"`
	CodeSource  *DeploymentCodeSource  `json:"codeSource"`
}

type DeploymentImageSource struct {
	Enabled      bool     `json:"enabled"`
	Name         string   `json:"name"`
	RegistryAuth ObjectID `json:"registryAuth,omitzero"`
}

type DeploymentCodeSource struct {
	Enabled        bool           `json:"enabled"`
	BuildTool      base.BuildTool `json:"buildTool"`
	DockerfilePath string         `json:"dockerfilePath"` // for BuildToolDockerfile only
	ImageTag       string         `json:"imageTag"`
	RegistryAuth   ObjectID       `json:"registryAuth,omitzero"`
}

func (s *Setting) AsAppDeploymentSettings() (*AppDeploymentSettings, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*AppDeploymentSettings)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &AppDeploymentSettings{}
	if s.Data != "" && s.Type == base.SettingTypeAppDeployment {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsAppDeploymentSettings() *AppDeploymentSettings {
	return gofn.Must(s.AsAppDeploymentSettings())
}
