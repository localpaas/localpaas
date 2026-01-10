package appdto

import (
	"github.com/docker/docker/api/types/swarm"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/services/docker"
)

type GetAppDeploymentSettingsReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
}

func NewGetAppDeploymentSettingsReq() *GetAppDeploymentSettingsReq {
	return &GetAppDeploymentSettingsReq{}
}

func (req *GetAppDeploymentSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppDeploymentSettingsResp struct {
	Meta *basedto.BaseMeta       `json:"meta"`
	Data *DeploymentSettingsResp `json:"data"`
}

type DeploymentSettingsResp struct {
	ImageSource   *DeploymentImageSourceResp   `json:"imageSource,omitempty"`
	RepoSource    *DeploymentRepoSourceResp    `json:"repoSource,omitempty"`
	TarballSource *DeploymentTarballSourceResp `json:"tarballSource,omitempty"`

	Command               string `json:"command,omitempty"`
	WorkingDir            string `json:"workingDir,omitempty"`
	PreDeploymentCommand  string `json:"preDeploymentCommand,omitempty"`
	PostDeploymentCommand string `json:"postDeploymentCommand,omitempty"`

	UpdateVer int `json:"updateVer"`
}

type DeploymentImageSourceResp struct {
	Enabled      bool                     `json:"enabled"`
	Image        string                   `json:"image"`
	RegistryAuth *basedto.NamedObjectResp `json:"registryAuth"`
}

type DeploymentRepoSourceResp struct {
	Enabled        bool                     `json:"enabled"`
	BuildTool      base.BuildTool           `json:"buildTool"`
	RepoURL        string                   `json:"repoUrl"`
	RepoRef        string                   `json:"repoRef"` // can be branch name, tag...
	Credentials    *RepoCredentialsResp     `json:"credentials"`
	DockerfilePath string                   `json:"dockerfilePath"` // for BuildToolDockerfile only
	ImageTags      []string                 `json:"imageTags"`
	RegistryAuth   *basedto.NamedObjectResp `json:"registryAuth"`
}

type RepoCredentialsResp struct {
	ID   string           `json:"id"`
	Type base.SettingType `json:"type"`
}

type DeploymentTarballSourceResp struct {
	Enabled bool `json:"enabled"`
}

type AppDeploymentSettingsTransformInput struct {
	App                *entity.App
	DeploymentSettings *entity.Setting
	ServiceSpec        *swarm.ServiceSpec
}

func TransformDeploymentSettings(input *AppDeploymentSettingsTransformInput) (resp *DeploymentSettingsResp, err error) {
	resp = &DeploymentSettingsResp{}

	if input.ServiceSpec != nil && input.ServiceSpec.TaskTemplate.ContainerSpec != nil {
		resp.WorkingDir = input.ServiceSpec.TaskTemplate.ContainerSpec.Dir
		resp.Command = docker.ConvertFromServiceCommand(input.ServiceSpec.TaskTemplate.ContainerSpec.Command,
			input.ServiceSpec.TaskTemplate.ContainerSpec.Args)
	}

	if input.DeploymentSettings != nil {
		if err = copier.Copy(&resp, input.DeploymentSettings); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	return resp, nil
}
