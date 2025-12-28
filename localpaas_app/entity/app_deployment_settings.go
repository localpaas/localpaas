package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentAppDeploymentSettingsVersion = 1
)

type AppDeploymentSettings struct {
	ImageSource   *DeploymentImageSource   `json:"imageSource"`
	RepoSource    *DeploymentRepoSource    `json:"repoSource"`
	TarballSource *DeploymentTarballSource `json:"tarballSource"`
}

type DeploymentImageSource struct {
	Enabled      bool     `json:"enabled"`
	Image        string   `json:"image"`
	RegistryAuth ObjectID `json:"registryAuth,omitzero"`
}

type DeploymentRepoSource struct {
	Enabled        bool           `json:"enabled"`
	BuildTool      base.BuildTool `json:"buildTool"`
	DockerfilePath string         `json:"dockerfilePath"` // for BuildToolDockerfile only
	ImageTag       string         `json:"imageTag"`
	RegistryAuth   ObjectID       `json:"registryAuth,omitzero"`
}

type DeploymentTarballSource struct {
	Enabled bool `json:"enabled"`
}

func (s *Setting) AsAppDeploymentSettings() (*AppDeploymentSettings, error) {
	return parseSettingAs(s, base.SettingTypeAppDeployment,
		func() *AppDeploymentSettings { return &AppDeploymentSettings{} })
}

func (s *Setting) MustAsAppDeploymentSettings() *AppDeploymentSettings {
	return gofn.Must(s.AsAppDeploymentSettings())
}
