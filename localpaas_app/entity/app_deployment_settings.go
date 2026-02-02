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
	ActiveMethod  base.DeploymentMethod    `json:"activeMethod"`

	Command               *string `json:"command,omitempty"`
	WorkingDir            *string `json:"workingDir,omitempty"`
	PreDeploymentCommand  *string `json:"preDeploymentCommand,omitempty"`
	PostDeploymentCommand *string `json:"postDeploymentCommand,omitempty"`
}

type DeploymentImageSource struct {
	Image        string   `json:"image"`
	RegistryAuth ObjectID `json:"registryAuth,omitzero"`
}

type DeploymentRepoSource struct {
	BuildTool      base.BuildTool  `json:"buildTool"`
	RepoType       base.RepoType   `json:"repoType"`
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
}

func (s *AppDeploymentSettings) GetType() base.SettingType {
	return base.SettingTypeAppDeployment
}

func (s *AppDeploymentSettings) GetRefSettingIDs() []string {
	res := make([]string, 0, 5) //nolint
	res = append(res, s.GetInUseRegistryAuthIDs()...)
	res = append(res, s.GetInUseGitCredentialIDs()...)
	return res
}

func (s *AppDeploymentSettings) GetInUseRegistryAuthIDs() (res []string) {
	if s.ImageSource != nil && s.ImageSource.RegistryAuth.ID != "" {
		res = append(res, s.ImageSource.RegistryAuth.ID)
	}
	if s.RepoSource != nil && s.RepoSource.RegistryAuth.ID != "" {
		res = append(res, s.RepoSource.RegistryAuth.ID)
	}
	res = gofn.ToSet(res)
	return
}

func (s *AppDeploymentSettings) GetInUseGitCredentialIDs() (res []string) {
	if s.RepoSource != nil && s.RepoSource.Credentials.ID != "" {
		res = append(res, s.RepoSource.Credentials.ID)
	}
	res = gofn.ToSet(res)
	return
}

func (s *Setting) AsAppDeploymentSettings() (*AppDeploymentSettings, error) {
	return parseSettingAs(s, func() *AppDeploymentSettings { return &AppDeploymentSettings{} })
}

func (s *Setting) MustAsAppDeploymentSettings() *AppDeploymentSettings {
	return gofn.Must(s.AsAppDeploymentSettings())
}
