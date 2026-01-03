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

	PreDeployment  *PreDeployment  `json:"preDeployment"`
	PostDeployment *PostDeployment `json:"postDeployment"`
}

type DeploymentImageSource struct {
	Enabled      bool     `json:"enabled"`
	Image        string   `json:"image"`
	RegistryAuth ObjectID `json:"registryAuth,omitzero"`
}

type DeploymentRepoSource struct {
	Enabled        bool            `json:"enabled"`
	BuildTool      base.BuildTool  `json:"buildTool"`
	RepoURL        string          `json:"repoUrl"`
	RepoRef        string          `json:"repoRef"`              // can be branch name, tag...
	Credentials    RepoCredentials `json:"credentials,omitzero"` // contains setting id of github app/git token/ssh key
	DockerfilePath string          `json:"dockerfilePath"`       // for BuildToolDockerfile only
	ImageTags      []string        `json:"imageTags"`
	RegistryAuth   ObjectID        `json:"registryAuth,omitzero"`
}

type RepoCredentials struct {
	ID   string           `json:"id"`
	Type base.SettingType `json:"type"`
}

type DeploymentTarballSource struct {
	Enabled bool `json:"enabled"`
}

type PreDeployment struct {
	Cmd string `json:"cmd"`
}

type PostDeployment struct {
	Cmd string `json:"cmd"`
}

func (s *Setting) AsAppDeploymentSettings() (*AppDeploymentSettings, error) {
	return parseSettingAs(s, base.SettingTypeAppDeployment,
		func() *AppDeploymentSettings { return &AppDeploymentSettings{} })
}

func (s *Setting) MustAsAppDeploymentSettings() *AppDeploymentSettings {
	return gofn.Must(s.AsAppDeploymentSettings())
}
