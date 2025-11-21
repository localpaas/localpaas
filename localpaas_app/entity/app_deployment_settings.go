package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type AppDeploymentSettings struct {
	ImageSource *DeploymentImageSource `json:"imageSource"`
	CodeSource  *DeploymentCodeSource  `json:"codeSource"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
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

func (s *Setting) ParseAppDeploymentSettings() (*AppDeploymentSettings, error) {
	res := &AppDeploymentSettings{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeAppDeployment {
		return res, s.parseData(res)
	}
	return nil, nil
}
