package appdto

import (
	"github.com/docker/docker/api/types/swarm"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
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
	ActiveSource  base.DeploymentSource        `json:"activeSource"`

	Command               string `json:"command,omitempty"`
	WorkingDir            string `json:"workingDir,omitempty"`
	PreDeploymentCommand  string `json:"preDeploymentCommand,omitempty"`
	PostDeploymentCommand string `json:"postDeploymentCommand,omitempty"`

	UpdateVer int `json:"updateVer"`
}

type DeploymentImageSourceResp struct {
	Image        string                    `json:"image"`
	RegistryAuth *settings.BaseSettingResp `json:"registryAuth"`
}

type DeploymentRepoSourceResp struct {
	BuildTool      base.BuildTool            `json:"buildTool"`
	RepoURL        string                    `json:"repoUrl"`
	RepoRef        string                    `json:"repoRef"` // can be branch name, tag...
	Credentials    *settings.BaseSettingResp `json:"credentials"`
	DockerfilePath string                    `json:"dockerfilePath"` // for BuildToolDockerfile only
	ImageTags      []string                  `json:"imageTags"`
	RegistryAuth   *settings.BaseSettingResp `json:"registryAuth"`
}

type DeploymentTarballSourceResp struct {
	// TODO: implement this
}

type AppDeploymentSettingsTransformInput struct {
	App                *entity.App
	DeploymentSettings *entity.Setting
	ServiceSpec        *swarm.ServiceSpec
	RefSettingMap      map[string]*entity.Setting
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
		appDeploymentSettings := input.DeploymentSettings.MustAsAppDeploymentSettings()
		if err = copier.Copy(&resp, appDeploymentSettings); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	if resp.ImageSource != nil { //nolint:nestif
		if resp.ImageSource.RegistryAuth != nil && resp.ImageSource.RegistryAuth.ID != "" {
			settingResp, _ := settings.TransformSettingBase(input.RefSettingMap[resp.ImageSource.RegistryAuth.ID])
			if settingResp != nil {
				resp.ImageSource.RegistryAuth = settingResp
			} else {
				resp.ImageSource.RegistryAuth = nil
			}
		}
	}
	if resp.RepoSource != nil { //nolint:nestif
		if resp.RepoSource.Credentials != nil && resp.RepoSource.Credentials.ID != "" {
			settingResp, _ := settings.TransformSettingBase(input.RefSettingMap[resp.RepoSource.Credentials.ID])
			if settingResp != nil {
				resp.RepoSource.Credentials = settingResp
			} else {
				resp.RepoSource.Credentials = nil
			}
		}
		if resp.RepoSource.RegistryAuth != nil && resp.RepoSource.RegistryAuth.ID != "" {
			settingResp, _ := settings.TransformSettingBase(input.RefSettingMap[resp.RepoSource.RegistryAuth.ID])
			if settingResp != nil {
				resp.RepoSource.RegistryAuth = settingResp
			} else {
				resp.RepoSource.RegistryAuth = nil
			}
		}
	}

	return resp, nil
}
